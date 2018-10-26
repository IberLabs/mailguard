package main

/**
More information about emerson IMAP lib
https://github.com/emersion/go-imap/wiki
 */

import (
	"log"
	"os"
	in "mailguard/internal"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
)

const CONST_MAXUNREADMESSAGESPERCYCLE		= 10
const CONST_INBOXFOLDERNAME					= "INBOX"

/**
	Open a new IMAP connection and return the connection pointer
 */
func openIMAPConnection(config in.Config) (c *client.Client) {
	// Connect to server
	connection, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Println("Connected")

	return connection
}

/**
	Check incoming messages
	TODO: Get message bodies (only envelope right now): https://github.com/emersion/go-imap/wiki/Fetching-messages
 */
func incomingMail(config in.Config){

	// Open connection
	c := openIMAPConnection(cfg)

	// Don't forget to logout at the end of this method
	defer c.Logout()

	// IMAP authentication
	if err := c.Login(config.Auth.Username, config.Auth.Password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// Point to Inbox folder
	mbox := getIMAPFolder(c, CONST_INBOXFOLDERNAME)

	// Get last unread e-mails
	messages, seqset := getLastUnreadMessages(c, mbox, CONST_MAXUNREADMESSAGESPERCYCLE)

	// Feed the message list
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	// Evaluate messages and trigger actions if needed
	evalAndTriggerActions(config, messages)

	if err := <-done; err != nil {
		log.Fatal(err)
	}
}

/**
	Return a pointer to a determined mailbox folder
 */
func getIMAPFolder(c * client.Client, foldername string) imap.MailboxStatus {
	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func () {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select folder
	mbox, err := c.Select(foldername, false)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Println("Flags for " + foldername + ":", mbox.Flags)

	return *mbox
}

/**
	Extract a chan struct with last N unread messages.
 	@limit means the max number of unread messages each time this method is invoked.
 */
func getLastUnreadMessages(c * client.Client, mbox imap.MailboxStatus, limit uint32) (chan*(imap.Message), *imap.SeqSet)  {
	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > limit {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - limit
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)

	return messages, seqset
}

