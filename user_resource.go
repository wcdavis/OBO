package main

import (
	"github.com/emicklei/go-restful"
	//"log"
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

	ws.Route(ws.GET("/{userId}").To(u.findUser).
		Doc("find a user").
		Operation("findUser").
		Param(ws.PathParameter("userId", "identifier of the user").DataType("string")).
		Writes(User{})) // on the response

	ws.Route(ws.PUT("/{userId}").To(u.updateUser).
		Doc("update a user").
		Operation("updateUser").
		Param(ws.PathParameter("userId", "identifier of the user").DataType("string")).
		ReturnsError(409, "duplicate userId", nil).
		Reads(User{})) // from the request

	ws.Route(ws.POST("").To(u.createUser).
		Doc("create a user").
		Operation("createUser").
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{userId}").To(u.removeUser).
		Doc("delete a user").
		Operation("deleteUser").
		Param(ws.PathParameter("userId", "identifier of the user").DataType("string")).
		Writes(User{}))

	container.Add(ws)
}

//--------------------------------------------------------------------//
// Request Functions

func (u *UserResource) findUser(request *restful.Request, response *restful.Response) {
	id, success := u.checkUserId(request, response)
	if !success {
		return
	}
	usr := u.storage.GetUser(id)
	response.WriteEntity(usr)
}

func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	usr, success := u.checkUser(request, response)
	if !success {
		return
	}
	usr.Id = u.generateUserId()
	u.storage.InsertUser(usr)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(usr)
}

func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
	id, success1 := u.checkUserId(request, response)
	usr, success2 := u.checkUser(request, response)
	if !success1 || !success2 {
		return
	}

	usr.Id = id // make sure the id is consistent

	u.storage.UpdateUser(usr)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(usr)
}

func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	id, success := u.checkUserId(request, response)
	if !success {
		return
	}

	usr := u.storage.DeleteUser(id)
	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(usr)
}

//--------------------------------------------------------------------//
// Utility Functions

func (u *UserResource) checkUserId(request *restful.Request, response *restful.Response) (int, bool) {
	success := true

	id, err := strconv.Atoi(request.PathParameter("userId"))
	if err != nil {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Malformed userId.")
	}
	if !u.storage.ExistsUser(id) {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "User not found.")
	}

	return id, success
}

func (u *UserResource) checkUser(request *restful.Request, response *restful.Response) (User, bool) {
	success := true

	usr := new(User)
	err := request.ReadEntity(usr)

	if err != nil {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed user.")
	}

	return *usr, success
}

func (u *UserResource) generateUserId() int {
	return u.storage.Length() + 1
}
