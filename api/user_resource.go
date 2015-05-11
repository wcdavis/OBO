package api

import (
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"

	. "github.com/PrincetonOBO/OBOBackend/item"
	. "github.com/PrincetonOBO/OBOBackend/user"

	"github.com/PrincetonOBO/OBOBackend/validate"
)

//type User user.User

type UserResource struct {
	storage     *UserStorage
	itemStorage *ItemStorage
	validator   *validate.Validator
}

func NewUserResource(db *mgo.Database) *UserResource {
	ur := new(UserResource)
	ur.storage = NewUserStorage(db)
	ur.itemStorage = NewItemStorage(db)
	ur.validator = validate.NewValidator(db)
	return ur
}

// significant boilerplate for registration adapted from
// https://github.com/emicklei/go-restful/blob/master/examples/restful-user-resource.go
func (u UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/manage/users").
		Doc("Manage Users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{user_id}").
		Filter(u.validator.Authenticate).
		Filter(u.validator.CheckUserId).
		To(u.findUser).
		Doc("find a user").
		Operation("findUser").
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Writes(User{})) // on the response

	ws.Route(ws.PUT("/{user_id}").
		Filter(u.validator.Authenticate).
		Filter(u.validator.CheckUserId).
		Filter(u.validator.CheckUser).
		To(u.updateUser).
		Doc("update a user").
		Operation("updateUser").
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Returns(409, "duplicate userId", nil).
		Reads(User{})) // from the request

	ws.Route(ws.POST("").
		Filter(u.validator.CheckUser).
		To(u.createUser).
		Doc("create a user").
		Operation("createUser").
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{user_id}").
		Filter(u.validator.Authenticate).
		Filter(u.validator.CheckUserId).
		To(u.removeUser).
		Doc("delete a user").
		Operation("deleteUser").
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Writes(User{}))

	ws.Route(ws.GET("/{user_id}/offers").
		Filter(u.validator.Authenticate).
		Filter(u.validator.CheckUserId).
		To(u.getActiveOffers).
		Doc("gets a user's active offers").
		Operation("getOffers").
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Writes([]ItemPresenter{}))

	container.Add(ws)
}

//--------------------------------------------------------------------//
// Request Functions

func (u *UserResource) findUser(request *restful.Request, response *restful.Response) {
	id := bson.ObjectIdHex(request.PathParameter("user_id"))
	usr := u.storage.GetUser(id)
	response.WriteEntity(usr)
}

func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	request.ReadEntity(usr)

	_, usr.Id = u.storage.InsertUser(*usr)
	usr.Authentication = u.validator.CreateAuthenticatedToken(*usr)

	u.storage.UpdateUser(*usr)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(usr)
}

func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
	id := bson.ObjectIdHex(request.PathParameter("user_id"))
	usr := new(User)
	request.ReadEntity(usr)

	usr.Id = id // make sure the id is consistent

	u.storage.UpdateUser(*usr)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(usr)
}

func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	id := bson.ObjectIdHex(request.PathParameter("user_id"))

	usr := u.storage.DeleteUser(id)
	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(usr)
}

func (u *UserResource) getActiveOffers(request *restful.Request, response *restful.Response) {
	id := bson.ObjectIdHex(request.PathParameter("user_id"))
	items := u.itemStorage.GetItemsByOffer(id)

	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(PresentWithOffer(*items, id))
}
