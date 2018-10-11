package main

import (
	"github.com/emersion/go-imap"
	"log"
	"github.com/emirpasic/gods/utils"
	"fmt"
)

/**
Receive e-mails in txt format and evaluate parameters.
Output will be a list of determined behaviour.
 */
func evalAndTriggerActions(config Config, messages chan * imap.Message)  {

	// Walk around message list
	log.Println("Last " + utils.ToString(CONST_MAXUNREADMESSAGESPERCYCLE) +  " messages:")
	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
		fmt.Println(msg.Envelope.Date)

		evalRuleAndTriggerAction(msg)
		
		fmt.Println()
	}
}

/**
	Evaluate a message and generate an action
 */
func evalRuleAndTriggerAction(msg *imap.Message) {

	// Rule 1 : Send e-mail to sender
	if(msg.Envelope.Subject == "hola"){
		fmt.Println("Rule 1 activated by message: " + msg.Envelope.Subject)

		// Send emails
		sendMail(	cfg, cfg.Account.Sender, (msg.Envelope.Sender[0].MailboxName + "@" + msg.Envelope.Sender[0].HostName),
			"Que tal vamos", "Que tal vamos. \r\nque tal estamos")

	}else{
		fmt.Println("No rules matched")
	}



}