package util

import (
	"encoding/json"
	"os"

	"github.com/emicklei/go-restful/log"
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

func Log(s string) {
	log.Print(s)
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
			dir, _ := os.Getwd()
			log.Print("error parsing " + os.Args[1] + " - using default values")
			log.Print("pwd = " + dir)
			configuration = defaultConfig
		}
	}

	configuration.WebURL = configuration.AppBaseURL + ":" + configuration.Port
	return configuration
}
