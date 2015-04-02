package util

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	DbName         string
	AppBaseURL     string
	DbBaseURL      string
	Port           string
	ApiPath        string
	WebURL         string
	SwaggerPath    string
	SwaggerBaseURL string
}

func Logerr(err error) {
	if err != nil {
		log.Print(err)
	}
}

func GetConfig() Configuration {
	defaultConfig := Configuration{DbName: "obo",
		AppBaseURL: "localhost", DbBaseURL: "localhost", Port: "4000", ApiPath: "/apidocs.json",
		SwaggerPath: "/apidocs/", SwaggerBaseURL: "/swagger/dist"}

	configuration := defaultConfig

	if len(os.Args) > 1 {
		file, _ := os.Open(os.Args[1])
		decoder := json.NewDecoder(file)
		configuration = Configuration{}
		err := decoder.Decode(&configuration)
		if err != nil {
			log.Print("error parsing conf.json - using default values")
			configuration = defaultConfig
		}
	}

	configuration.WebURL = configuration.AppBaseURL + ":" + configuration.Port
	return configuration
}
