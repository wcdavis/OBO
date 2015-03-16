package main

import (
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"strconv"
)

type UserResource struct {
	storage *UserStorage
}

// significant boilerplate for registration adapted from
// https://github.com/emicklei/go-restful/blob/master/examples/restful-user-resource.go
func (u UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Doc("Manage Users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		// docs
		Doc("get a user").
		Operation("findUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(User{})) // on the response
	/*
		ws.Route(ws.PUT("/{user-id}").To(u.updateUser).
			// docs
			Doc("update a user").
			Operation("updateUser").
			Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
			ReturnsError(409, "duplicate user-id", nil).
			Reads(User{})) // from the request
	*/
	ws.Route(ws.POST("").To(u.createUser).
		Doc("create a user").
		Operation("createUser").
		Reads(User{})) // from the request
	/*
		ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
			Doc("delete a user").
			Operation("removeUser").
			Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))
	*/
	container.Add(ws)
}

func (u *UserResource) findUser(request *restful.Request, response *restful.Response) {
	id, err := strconv.Atoi(request.PathParameter("user-id"))
	if (err != nil) || (!u.storage.ExistsUser(id)) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}
	usr := u.storage.GetUser(id)
	response.WriteEntity(usr)
}

func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(usr)

	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf(usr.FirstName)
	usr.Id = u.storage.Length() + 1 // simple id generation !!! CHANGE
	u.storage.InsertUser(*usr)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(usr)
}
