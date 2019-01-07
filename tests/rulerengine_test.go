package main

import (
	"testing"
	in "mailguard/internal"
	"log"
	"os"
	"time"
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

func TestSum(t *testing.T) {
	total := 10
	if total != 10 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	}
}

func TestRules(t *testing.T) {

	// Initialize an empty slice of DataUnits
	dataUnitList := make([]in.DataUnit, 2)
	currTime, _ := time.Parse("2006-01-02 00:00:00", "2019-01-19 08:00:00");

	dataUnitList[0] = in.DataUnit{
		Date: 		currTime,
		Subject: 	"This is a test",
		From: 		"test@test.com",
		To:			"test2@test2.com",
	}

	dataUnitList[1] = in.DataUnit{
		Date: 		currTime,
		Subject: 	"This is another test",
		From: 		"test@test.com",
		To:			"test2@test2.com",
	}

	// Populate test rules
	rules := []string{
		"subject == 'hello'",
		"time() > '19:00:00' || (dayOfWeek() == 0 || dayOfWeek() == 6)",
	}

	in.EvalRulesAndTriggerActions(&cfg, rules, &dataUnitList)
}