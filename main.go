package main

import (
	"log"

	"github.com/PrincetonOBO/OBOBackend/user"
	"github.com/PrincetonOBO/OBOBackend/util"

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
	log.Print("Establishing connection with Mongo at " + configuration.DbBaseURL + "...")
	session, err := mgo.Dial(configuration.DbBaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	log.Print("Connection established!")

	database := session.DB(configuration.DbName)

	// create resources
<<<<<<< HEAD
	userResource := UserResource{NewUserStorage(database)}
	itemResource := ItemResource{NewItemStorage(database)}

	// register services
	userResource.Register(wsContainer)
	itemResource.Register(wsContainer)
=======
	resources := createResources(database)

	// register services
	for _, resource := range resources {
		resource.Register(wsContainer)
	}
>>>>>>> 451c3be828f6c772df6afad60b7fb60d7cd16ffa

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
