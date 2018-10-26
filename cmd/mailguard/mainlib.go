package main

import (
	"os"
	"fmt"
	"log"
	"encoding/json"
	"path/filepath"
	in "mailguard/internal"
)

func loadConfiguration(file string) in.Config {
	var configData in.Config
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
