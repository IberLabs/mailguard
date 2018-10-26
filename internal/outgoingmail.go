package internal

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"strings"
	"log"
)

/**
	Send e-mail messages
*/
func sendMail(config Config, from string, to string, subject string, body string) {
	// Send message
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	auth := sasl.NewPlainClient("", config.Auth.Username, config.Auth.Password)

	smtpTo := []string{to}
	msg := strings.NewReader("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\n")
	smtpErr := smtp.SendMail(config.Smtp.Host + ":" + config.Smtp.Port, auth, config.Account.Sender, smtpTo, msg)
	if smtpErr != nil {
		log.Fatal(smtpErr)
	}
}
