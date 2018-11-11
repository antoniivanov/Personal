package main

import (
	"io/ioutil"
	"log"
	"os"
)

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func readFileToString(file string) string {
	readBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln("Cannot read answer file", file)
	}
	return string(readBytes)
}
