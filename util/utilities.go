package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	DbName         string
	BaseURL        string
	Port           string
	ApiPath        string
	WebURL         string
	SwaggerPath    string
	SwaggerBaseURL string
}

type Configuration1 struct {
	Users  []string
	Groups []string
}

func logerr(err error) {
	if err != nil {
		log.Print(err)
	}
}

func getConfig() Configuration {
	defaultConfig := Configuration{DbName: "obo",
		BaseURL: "localhost", Port: "80", ApiPath: "/apidocs.json",
		SwaggerPath: "/apidocs/", SwaggerBaseURL: "/swagger/dist"}

	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Print("error parsing conf.json - using default values")
		configuration = defaultConfig
	}
	configuration.WebURL = configuration.BaseURL + ":" + configuration.Port
	return configuration
}
