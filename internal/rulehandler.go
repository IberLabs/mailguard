package internal

import (
	"fmt"
	"log"
	"github.com/emirpasic/gods/utils"
	"strings"
	"gopkg.in/Knetic/govaluate.v2"
)

const CONST_TAG_EVAL	= "eval:"

/**
Receive e-mails in txt format and evaluate parameters.
Output will be a list of determined behaviour.
 */
func EvalRulesAndTriggerActions(config Config, rules []string, dataUnitList * [] DataUnit, maxUnreadMessagesPerCycle int)  {

	// Walk around message list
	log.Println("Last " + utils.ToString(maxUnreadMessagesPerCycle) +  " messages:")
	for element := range *dataUnitList {
		evalRuleAndTriggerAction(&(*dataUnitList)[element], config, rules)
	}
}

/**
	Evaluate a rule against a particular dataUnits an trigger an action.
	Both rules and actions are loaded from yml file.
 */
func evalRuleAndTriggerAction(dataUnit * DataUnit, config Config, rules []string) {

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


		// Add message fields to evaluation parameters
		parameters["date"] = dataUnit.Date.String()
		parameters["from"] = dataUnit.From
		parameters["to"] = dataUnit.To
		parameters["subject"] = dataUnit.Subject


		/** TODO: Add functions evaluation */
		result, err := expression.Evaluate(parameters);

		if err == nil && result == true {
			rulesMatched = true

			//Action temporary hardcoded
			fmt.Println("Rule 1 activated by message: " + dataUnit.Subject)
			fmt.Println("Sending message (Hardcoded for now)")

			// Send emails
			sendMail(	config, config.Account.Sender, (dataUnit.From),
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