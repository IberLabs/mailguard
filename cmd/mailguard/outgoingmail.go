package main

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"strings"
	"log"
)

/**
	Send e-mail messages
*/
func outgoingMail(config Config) {
	// Send message
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	auth := sasl.NewPlainClient("", config.Auth.Username, config.Auth.Password)
	smtpTo := []string{"manuelbcd@gmail.com"}
	msg := strings.NewReader("To: manuelbcd@gmail.com\r\n" +
		"Subject: Una pruebecilla\r\n" +
		"\r\n" +
		"Esto es una prueba.\r\n")
	smtpErr := smtp.SendMail("smtp.gmail.com:587", auth, config.Account.Sender, smtpTo, msg)
	if smtpErr != nil {
		log.Fatal(smtpErr)
	}
}
