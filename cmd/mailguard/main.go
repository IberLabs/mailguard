package main

import (
	"log"
	"os"
)

const configFilename = 	"config.json"

var appDir 			string			// Configuration var for app folder
var cfg 			Config			// Main configuration var


func init() {
	appDir 		  	:= getAppDir()
	if appDir == "" {
		log.Println("Error retrieving app dir")
		os.Exit(0)
	}

	cfg 	= loadConfiguration(appDir + "/" + configFilename)
	println(cfg.Imap.Host)
	println(cfg.Auth.Username)
	if cfg.Auth.Username == "" {
		log.Println("Error: username can't be empty")
		os.Exit(0)
	}
	if cfg.Auth.Password == "" {
		log.Println("Error: password can't be empty")
		os.Exit(0)
	}
}

func main() {
	log.Println("Connecting to server...")

	// Check mailbox
	incomingMail(cfg)



	// Send emails
	outgoingMail(cfg)

	log.Println("Done!")
}