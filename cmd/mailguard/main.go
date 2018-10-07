package main

import (
	"log"
	"strings"

	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"os"
	"fmt"
	"encoding/json"
	"path/filepath"
)

type Config struct {
	Smtp 	struct {
		Host     	string `json:"host"`
		Port	 	string `json:"port"`
	} `json:"struct"`
	Imap 	struct {
		Host     	string `json:"host"`
		Port 		string `json:"port"`
	} `json:"imap"`
	Auth 	struct {
		Username	string `json:"username"`
		Password 	string `json:"password"`
	} `json:"auth"`
	Account struct{
		Sender		string `json:"sender"`
	}
}

const configFilename = 	"config.json"

var appDir 			string
var configuration 	Config

func init() {
	appDir 		  	:= getAppDir()
	if appDir == "" {
		log.Println("Error retrieving app dir")
		os.Exit(0)
	}

	configuration 	= loadConfiguration(appDir + "/" + configFilename)
	println(configuration.Imap.Host)
	println(configuration.Auth.Username)
	if configuration.Auth.Username == "" {
		log.Println("Error: username can't be empty")
		os.Exit(0)
	}
	if configuration.Auth.Password == "" {
		log.Println("Error: password can't be empty")
		os.Exit(0)
	}
}

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(configuration.Auth.Username, configuration.Auth.Password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

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

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 3 {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - 3
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	log.Println("Last 4 messages:")
	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}


	// Send message
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	auth := sasl.NewPlainClient("", configuration.Auth.Username, configuration.Auth.Password)
	smtpTo := []string{"manuelbcd@gmail.com"}
	msg := strings.NewReader("To: manuelbcd@gmail.com\r\n" +
		"Subject: Una pruebecilla\r\n" +
		"\r\n" +
		"Esto es una prueba.\r\n")
	smtpErr := smtp.SendMail("smtp.gmail.com:587", auth, configuration.Account.Sender, smtpTo, msg)
	if smtpErr != nil {
		log.Fatal(smtpErr)
	}

	log.Println("Done!")
}

func loadConfiguration(file string) Config {
	var configData Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&configData)

	return configData
}

func getAppDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return dir
}