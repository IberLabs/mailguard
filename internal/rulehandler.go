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
	Evaluate a message and generate an action
 */
func evalRuleAndTriggerAction(msg *imap.Message, config Config, rules []string) {

	// Walk around rules
	for _, rule := range rules {


		if !strings.Contains( rule, CONST_TAG_EVAL){
			continue
		}

		i := strings.Index(rule, CONST_TAG_EVAL)
		// Trim any empty char or tabulation
		rawRule := strings.Trim(rule[i+len(CONST_TAG_EVAL):], "	")
		rawRule = strings.Trim(rawRule, "	")
		rawRule = strings.Trim(rawRule, " ")

		fmt.Println("Rule detected")
		fmt.Println("'" + rawRule + "'")

		expression, err := govaluate.NewEvaluableExpression(rawRule);
		parameters := make(map[string]interface{}, 8)
		parameters["subject"] = "hello"
		result, err := expression.Evaluate(parameters);
		// result is now set to "false", the bool value
		_ = result	// temp
		_ = err		// temp

	}
	return //temp


	// TEMPORARY HARD-CODED---------------------------------------


	// Rule 1 : Send e-mail to sender
	if(msg.Envelope.Subject == "hola"){
		fmt.Println("Rule 1 activated by message: " + msg.Envelope.Subject)

		// Send emails
		sendMail(	config, config.Account.Sender, (msg.Envelope.Sender[0].MailboxName + "@" + msg.Envelope.Sender[0].HostName),
			"Que tal vamos", "Que tal vamos. \r\nque tal estamos")

	}else{
		fmt.Println("No rules matched")
	}

}