package internal

import (
	"os"
	"fmt"
	"log"
	"encoding/json"
	"bufio"
	"path/filepath"
)

func LoadConfiguration(file string) Config {
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

/**
	Obtain app dir (app executable folder with the absolute path)
 */
func GetAppDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

/**
	Load rules file
 */
func ReadFileLines(filePath string) (bool, []string) {
	readResult := true
	var lines []string
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		readResult = false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		//fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		readResult = false
	}

	return readResult, lines
}
