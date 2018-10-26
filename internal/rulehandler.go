package internal

import (
	"fmt"
	"log"
	"github.com/emirpasic/gods/utils"
	"github.com/emersion/go-imap"
	"gopkg.in/Knetic/govaluate.v2"
)

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


	/**
	TODO: Evaluate rules, create functions, contexts, logs.
	 */
	expression, err := govaluate.NewEvaluableExpression("foo > 0");

	if err != nil {
		log.Fatal(err)
	}

	parameters := make(map[string]interface{}, 8)
	parameters["foo"] = -1;

	result, err := expression.Evaluate(parameters);

	result = result

	if err != nil {
		log.Fatal(err)
	}



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