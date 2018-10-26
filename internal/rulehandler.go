package internal

import (
	"fmt"
	"log"
	"github.com/emirpasic/gods/utils"
	"github.com/emersion/go-imap"
)

/**
Receive e-mails in txt format and evaluate parameters.
Output will be a list of determined behaviour.
 */
func EvalRulesAndTriggerActions(config Config, messages chan * imap.Message, maxUnreadMessagesPerCycle int)  {

	// Walk around message list
	log.Println("Last " + utils.ToString(maxUnreadMessagesPerCycle) +  " messages:")
	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
		fmt.Println(msg.Envelope.Date)

		evalRuleAndTriggerAction(msg, config)
		
		fmt.Println()
	}
}

/**
	Evaluate a message and generate an action
 */
func evalRuleAndTriggerAction(msg *imap.Message, config Config) {

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