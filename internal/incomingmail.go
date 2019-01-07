package internal

/**
More information about emerson IMAP lib
https://github.com/emersion/go-imap/wiki
 */

import (
	"log"
	"os"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"fmt"
	"io"
	"io/ioutil"
)

const CONST_MAXUNREADMESSAGESPERCYCLE		= 10
const CONST_INBOXFOLDERNAME					= "INBOX"

/**
	Open a new IMAP connection and return the connection pointer
 */
func openIMAPConnection(config * Config) (c *client.Client) {
	// Connect to server
	connection, err := client.DialTLS(config.Imap.Host + ":" + config.Imap.Port, nil)
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
func IncomingMail(config * Config, rules []string){

	// Open connection
	c := openIMAPConnection(config)

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

	// Read full msg body
	section := &imap.BodySectionName{}

	// Feed the message list
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
	}()

	// Initialize an empty slice of DataUnits
	dataUnitList := make([]DataUnit, cap(messages)+1)

	// Get a list of DataUnits from e-mail messages
	convertEmailsToDataUnits(messages, &dataUnitList)

	// Evaluate messages and trigger actions if needed
	EvalRulesAndTriggerActions(config, rules, &dataUnitList)

	if err := <-done; err != nil {
		log.Fatal(err)
	}
}

/**
	Convert and populate a DataUnit list from E-mail list.
 */
func convertEmailsToDataUnits( messages chan * imap.Message, dataList * [] DataUnit) {

	var counter int = 0

	// Walk around message list
	for msg := range messages {

		// Get full body message
		section := &imap.BodySectionName{}

		r := msg.GetBody(section)
		if r == nil {
			log.Fatal("Server didn't return message body")
		}

		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Fatal(err)
		}

		// Add message fields to evaluation parameters
		header := mr.Header
		if date, err := header.Date(); err == nil {
			(*dataList)[counter].Date = date
		}

		if from, err := header.AddressList("From"); err == nil {
			(*dataList)[counter].From = from[0].Address
		}

		if to, err := header.AddressList("To"); err == nil {
			(*dataList)[counter].To = to[0].Address
		}

		if subject, err := header.Subject(); err == nil {
			(*dataList)[counter].Subject = subject
		}

		// Process each message's part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch p.Header.(type) {
			case mail.TextHeader:
				// This is the message's text (can be plain-text or HTML)
				b, _ := ioutil.ReadAll(p.Body)
				(*dataList)[counter].Body = string(b)
			case mail.AttachmentHeader:
				// This is an attachment
				//filename, _ := h.Filename()
				//log.Println("Got attachment: %v", filename)
				// TODO : What to do with attachments?
			}
		}

		fmt.Println((*dataList)[counter].From + " : " + (*dataList)[counter].Subject )
		fmt.Println("---------------------")

		counter++
	}

	return

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

	messages := make(chan *imap.Message, limit)

	return messages, seqset
}

