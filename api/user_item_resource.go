package api

import (
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"
	"time"

	. "github.com/PrincetonOBO/OBOBackend/image"
	. "github.com/PrincetonOBO/OBOBackend/item"
	. "github.com/PrincetonOBO/OBOBackend/user"

	"github.com/PrincetonOBO/OBOBackend/validate"
)

type UserItemResource struct {
	storage      *ItemStorage
	userStorage  *UserStorage
	imageStorage *ImageStorage
	validator    *validate.Validator
}

func NewUserItemResource(db *mgo.Database) *UserItemResource {
	uir := new(UserItemResource)
	uir.storage = NewItemStorage(db)
	uir.userStorage = NewUserStorage(db)
	uir.imageStorage = NewImageStorage(db)
	uir.validator = validate.NewValidator(db)
	return uir
}

// significant boilerplate for registration adapted from go-restful package
func (i UserItemResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/users/{user_id}/items").
		Filter(i.validator.Authenticate).
		Filter(i.validator.CheckUserId).
		Doc("Manage a User's Items").
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/").
		To(i.getUserItems).
		Doc("get a user's items").
		Operation("getUserItems").
		Writes([]Item{}))

	ws.Route(ws.POST("/").
		Filter(i.validator.CheckItem).
		To(i.createItem).
		Doc("create an item").
		Operation("createItem").
		Reads(ItemPresenter{})) // from the request

	ws.Route(ws.POST("/{item_id}/pic").
		Filter(i.validator.CheckImage).
		Filter(i.validator.CheckItemId).
		Filter(i.validator.CheckItemOwnership).
		To(i.addPicture).
		Doc("attach an item to the picture").
		Operation("attachPic").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Reads(ImagePresenter{})) // from the request

	ws.Route(ws.GET("/{item_id}").
		Filter(i.validator.CheckItemId).
		Filter(i.validator.CheckItemOwnership).
		To(i.findItem).
		Doc("find an item").
		Operation("findUserItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Writes(Item{})) // on the response

	ws.Route(ws.GET("/{item_id}/offer/{offer_id}").
		Filter(i.validator.CheckOfferId).
		Filter(i.validator.CheckItemId).
		To(i.acceptOffer).
		Doc("accept an offer").
		Operation("acceptOffer").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Param(ws.PathParameter("offer_id", "user-identifier of the offer").DataType("string")).
		Writes(OfferPresenter{})) // on the response

	ws.Route(ws.PUT("/{item_id}").
		Filter(i.validator.CheckItemId).
		Filter(i.validator.CheckItem).
		Filter(i.validator.CheckItemOwnership).
		To(i.updateItem).
		Doc("update an item").
		Operation("updateItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Reads(ItemPresenter{})) // from the request

	ws.Route(ws.DELETE("/{item_id}").
		Filter(i.validator.CheckItemId).
		To(i.removeItem).
		Doc("delete an item").
		Operation("deleteItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Writes(Item{}))

	container.Add(ws)
}

//--------------------------------------------------------------------//
// Request Functions

func (i *UserItemResource) findItem(request *restful.Request, response *restful.Response) {
	id := bson.ObjectIdHex(request.PathParameter("item_id"))

	item := i.storage.GetItem(id)
	response.WriteEntity(item)
}

func (i *UserItemResource) createItem(request *restful.Request, response *restful.Response) {
	itemPres := new(ItemPresenter)
	request.ReadEntity(itemPres)
	item := itemPres.ToItem()

	uid := bson.ObjectIdHex(request.PathParameter("user_id"))

	item.User_Id = uid
	item.Location.Type = "Point"
	item.Time = time.Now().Unix()

	_, item.Id = i.storage.InsertItem(item)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(item)
}

func (i *UserItemResource) addPicture(request *restful.Request, response *restful.Response) {
	imagePres := new(ImagePresenter)
	request.ReadEntity(imagePres)
	im := imagePres.ToImage()

	id := bson.ObjectIdHex(request.PathParameter("item_id"))
	im.Item_Id = id

	i.imageStorage.InsertImage(im)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(im.ToPresenter(THUMB))

}

func (i *UserItemResource) updateItem(request *restful.Request, response *restful.Response) {
	itemPres := new(ItemPresenter)
	request.ReadEntity(itemPres)
	item := itemPres.ToItem()

	id := bson.ObjectIdHex(request.PathParameter("item_id"))
	item.Id = id // make sure the id is consistent

	i.storage.UpdateItem(item)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(item)
}

func (i *UserItemResource) acceptOffer(request *restful.Request, response *restful.Response) {
	id := bson.ObjectIdHex(request.PathParameter("item_id"))
	offId := bson.ObjectIdHex(request.PathParameter("offer_id"))

	storedItem := i.storage.GetItem(id)

	for _, off := range storedItem.Offers {
		if off.User_Id == offId {
			// here's where we would put some sort of notification
			storedItem.Sold = true
			i.storage.UpdateItem(*storedItem)
			response.WriteHeader(http.StatusAccepted)
			response.WriteEntity(off)
			return
		}
	}

	// this is a bit hacky, maybe some other time i'll put it elsewhere
	// i don't like to have this sort of error in the body, it should be
	// in the body.
	response.WriteErrorString(http.StatusNotFound, "No such offer exists")
	return
}

func (i *UserItemResource) removeItem(request *restful.Request, response *restful.Response) {
	id := bson.ObjectIdHex(request.PathParameter("item_id"))

	item := i.storage.DeleteItem(id)
	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(item)
}

func (i *UserItemResource) getUserItems(request *restful.Request, response *restful.Response) {
	uid := bson.ObjectIdHex(request.PathParameter("user_id"))

	items := i.storage.GetItemsByUserId(uid)

	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(items)
}
