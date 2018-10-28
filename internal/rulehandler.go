package internal

import (
	"fmt"
	"log"
	"github.com/emirpasic/gods/utils"
	"github.com/emersion/go-imap"
	"strings"
	"gopkg.in/Knetic/govaluate.v2"
)

const CONST_TAG_EVAL	= "eval:"

/**
Receive e-mails in txt format and evaluate parameters.
Output will be a list of determined behaviour.
 */
func EvalRulesAndTriggerActions(config Config, rules []string, messages chan * imap.Message, maxUnreadMessagesPerCycle int)  {

	// Walk around message list
	log.Println("Last " + utils.ToString(maxUnreadMessagesPerCycle) +  " messages:")
	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
		fmt.Println(msg.Envelope.Date)

		evalRuleAndTriggerAction(msg, config, rules)
		fmt.Println()
	}
}

/**
	Evaluate a rule against messages an trigger an action.
	Both rules and actions are loaded from yml file.
 */
func evalRuleAndTriggerAction(msg *imap.Message, config Config, rules []string) {

	var rulesMatched = false

	// Walk around rules
	for _, rule := range rules {

		if !strings.Contains( rule, CONST_TAG_EVAL){
			continue
		}

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in after rule fail", r)
			}
		}()

		i := strings.Index(rule, CONST_TAG_EVAL)
		// Trim any empty char or tabulation
		rawRule := strings.Trim(rule[i+len(CONST_TAG_EVAL):], "	")
		rawRule = strings.Trim(rawRule, " ")

		fmt.Println("Rule detected")
		fmt.Println("'" + rawRule + "'")

		expression, err := govaluate.NewEvaluableExpression(rawRule);
		parameters := make(map[string]interface{}, 8)
		parameters["subject"] = msg.Envelope.Subject;
		parameters["sender"] = msg.Envelope.Sender;

		/** TODO: Add functions evaluation */

		result, err := expression.Evaluate(parameters);

		if err == nil && result == true {
			rulesMatched = true

			//Action temporary hardcoded
			fmt.Println("Rule 1 activated by message: " + msg.Envelope.Subject)
			fmt.Println("Sending message (Hardcoded for now)")

			// Send emails
			sendMail(	config, config.Account.Sender, (msg.Envelope.Sender[0].MailboxName + "@" + msg.Envelope.Sender[0].HostName),
				"Que tal vamos", "Que tal vamos. \r\nque tal estamos")
			break
		}


	}



	// TEMPORARY HARD-CODED---------------------------------------
	// Rule 1 : Send e-mail to sender
	/**
	if(msg.Envelope.Subject == "hola"){
		fmt.Println("Rule 1 activated by message: " + msg.Envelope.Subject)

		// Send emails
		sendMail(	config, config.Account.Sender, (msg.Envelope.Sender[0].MailboxName + "@" + msg.Envelope.Sender[0].HostName),
			"Que tal vamos", "Que tal vamos. \r\nque tal estamos")

	}else{
		fmt.Println("No rules matched")
	}
	*/

	if !rulesMatched {
		fmt.Println("No rules matched")
	}

}