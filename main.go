package main

import (
	"log"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"

	"gopkg.in/mgo.v2"
	"net/http"
	//"gopkg.in/mgo.v2/bson"
)

func main() {

	dbName := "obo"
	url := "localhost"

	// create a new web service container
	wsContainer := restful.NewContainer()

	// create database
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	database := session.DB(dbName)

	// create resources
	userResource := UserResource{NewUserStorage(database)}
	itemResource := ItemResource{NewItemStorage(database)}

	// register services
	userResource.Register(wsContainer)
	itemResource.Register(wsContainer)

	// configure swagger
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:4000",
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/Users/Ben/go/src/github.com/swagger/swagger-ui/dist"}
	swagger.RegisterSwaggerService(config, wsContainer)

	log.Printf("start listening on localhost:4000")
	server := &http.Server{Addr: ":4000", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())

}
