package user

import (
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"
)

type UserResource struct {
	storage *UserStorage
}

func NewUserResource(db *mgo.Database) *UserResource {
	ur := new(UserResource)
	ur.storage = newUserStorage(db)
	return ur
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

	ws.Route(ws.GET("/{user_id}").To(u.findUser).
		Doc("find a user").
		Operation("findUser").
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Writes(User{})) // on the response

	ws.Route(ws.PUT("/{user_id}").To(u.updateUser).
		Doc("update a user").
		Operation("updateUser").
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Returns(409, "duplicate userId", nil).
		Reads(User{})) // from the request

	ws.Route(ws.POST("").To(u.createUser).
		Doc("create a user").
		Operation("createUser").
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{user_id}").To(u.removeUser).
		Doc("delete a user").
		Operation("deleteUser").
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
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

	_, usr.Id = u.storage.InsertUser(usr)
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

func (u *UserResource) checkUserId(request *restful.Request, response *restful.Response) (bson.ObjectId, bool) {
	success := true
	idString := request.PathParameter("user_id")

	if !bson.IsObjectIdHex(idString) {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed user id.")
	} else if !u.storage.ExistsUser(bson.ObjectIdHex(idString)) {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "User not found.")
	}
	id := bson.ObjectIdHex(idString)

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
