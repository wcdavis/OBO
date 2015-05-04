package main

import (
	"github.com/PrincetonOBO/OBOBackend/api"
	//myim "github.com/PrincetonOBO/OBOBackend/image"
	"github.com/PrincetonOBO/OBOBackend/util"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
	"github.com/emicklei/go-restful/swagger"

	//"fmt"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	//"image"
	//"image/jpeg"
	"net/http"
	//"os"
)

type Resource interface {
	Register(*restful.Container)
}

func main() {
	///// START TEST
	/*
		image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

		imgfile, err := os.Open("./test2.jpg")

		if err != nil {
			fmt.Println("img.jpg file not found!")
			os.Exit(1)
		}

		defer imgfile.Close()

		img, _, err := image.Decode(imgfile)

		thisIm := myim.NewImage(img, bson.NewObjectId())
		fmt.Println(thisIm.Thumbnail)

		fmt.Println(img.At(10, 10))

		bounds := img.Bounds()

		fmt.Println(bounds)

		canvas := image.NewAlpha(bounds)

		// is this image opaque
		op := canvas.Opaque()

		fmt.Println(op)

		/////// END
	*/
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
	log.Print(server.ListenAndServe())

}

func createResources(database *mgo.Database) []Resource {
	var resources []Resource
	resources = append(resources, api.NewUserResource(database))
	resources = append(resources, api.NewUserItemResource(database))
	resources = append(resources, api.NewItemResource(database))

	return resources
}
