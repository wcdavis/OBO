package main

import (
	"log"

	"github.com/beisner/OBO/user"
	"github.com/beisner/OBO/util"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"

	"gopkg.in/mgo.v2"
	"net/http"
)

type Resource interface {
	Register(*restful.Container)
}

func main() {

	// get configuration
	configuration := util.GetConfig()

	// create a new web service container
	wsContainer := restful.NewContainer()

	// create database
	session, err := mgo.Dial(configuration.BaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	database := session.DB(configuration.DbName)

	// create resources
	resources := createResources(database)

	// register services
	for _, resource := range resources {
		resource.Register(wsContainer)
	}

	// configure swagger
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://" + configuration.WebURL,
		ApiPath:        configuration.ApiPath,

		// Specifiy where the UI is located
		SwaggerPath:     configuration.SwaggerPath,
		SwaggerFilePath: configuration.SwaggerBaseURL}
	swagger.RegisterSwaggerService(config, wsContainer)

	log.Printf("start listening on " + configuration.WebURL)
	server := &http.Server{Addr: ":" + configuration.Port, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())

}

func createResources(database *mgo.Database) []Resource {
	var resources []Resource
	resources = append(resources, user.NewUserResource(database))
	// resources[1] = ItemResource{}
	return resources
}
