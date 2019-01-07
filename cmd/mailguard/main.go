package main

import (
	"log"
	"os"
	in "mailguard/internal"
)

const configFilename = 	"config.json"
const rulesFilename  = "rules.yml"

var appDir 			string				// Configuration var for app folder
var cfg 			in.Config			// Main configuration var
var rules			[]string			// Rule list

func init() {
	appDir 		  	:= in.GetAppDir()
	if appDir == "" {
		log.Println("Error retrieving app dir")
		os.Exit(0)
	}

	cfg = in.LoadConfiguration(appDir + "/" + configFilename)
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

	_, rules = in.ReadFileLines(appDir + "/" + rulesFilename)
}

func main() {
	// Start
	log.Println("Connecting to server...")

	// Check mailbox
	in.IncomingMail(&cfg, rules)

	// End
	log.Println("Done!")
}