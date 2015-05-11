package api

import (
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"

	. "github.com/PrincetonOBO/OBOBackend/image"
	. "github.com/PrincetonOBO/OBOBackend/item"
	. "github.com/PrincetonOBO/OBOBackend/user"

	"github.com/PrincetonOBO/OBOBackend/util"
	"github.com/PrincetonOBO/OBOBackend/validate"

	"strconv"
)

type ItemResource struct {
	storage      *ItemStorage
	userStorage  *UserStorage
	imageStorage *ImageStorage
	validator    *validate.Validator
}

func NewItemResource(db *mgo.Database) *ItemResource {
	ir := new(ItemResource)
	ir.storage = NewItemStorage(db)
	ir.userStorage = NewUserStorage(db)
	ir.imageStorage = NewImageStorage(db)
	ir.validator = validate.NewValidator(db)
	return ir
}

// significant boilerplate for registration adapted from
// https://github.com/emicklei/go-restful/blob/master/examples/restful-user-resource.go
func (i ItemResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/items").
		Doc("Interact with Items").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/").
		Filter(i.validator.CheckFeedQuery).
		To(i.getFeed).
		Doc("get a feed of items").
		Operation("getFeed").
		Param(ws.QueryParameter("longitude", "longitude for query").DataType("float64")).
		Param(ws.QueryParameter("latitude", "longitude for query").DataType("float64")).
		Param(ws.QueryParameter("number", "number of entries to return").DataType("int")).
		Writes([]ItemListPresenter{}))

	ws.Route(ws.GET("/{item_id}").
		Filter(i.validator.CheckItemId).
		To(i.findItem).
		Doc("find an item").
		Operation("findItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Writes(ItemPresenter{})) // on the response

	ws.Route(ws.GET("/{item_id}/pic/{pic_id}").
		Filter(i.validator.CheckImageId).
		To(i.getPicture).
		Doc("get a picture's item").
		Operation("getPic").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Param(ws.PathParameter("pic_id", "identifier of the picture").DataType("string")).
		Writes(ImagePresenter{})) // from the request

	ws.Route(ws.POST("/{item_id}/offer/users/{user_id}").
		Filter(i.validator.Authenticate).
		Filter(i.validator.CheckItemId).
		Filter(i.validator.CheckUserId).
		Filter(i.validator.CheckOffer).
		To(i.postOffer).
		Doc("post an offer").
		Operation("newOffer").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		// may be replaced by user token
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Reads(OfferPresenter{})) // from the request

	ws.Route(ws.DELETE("/{item_id}/offer/users/{user_id}").
		Filter(i.validator.Authenticate).
		Filter(i.validator.CheckItemId).
		Filter(i.validator.CheckUserId).
		To(i.deleteOffer).
		Doc("delete an offer").
		Operation("deleteOffer").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		// may be replaced
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Writes(OfferPresenter{})) // from the request

	ws.Route(ws.GET("/{item_id}/report").
		Filter(i.validator.CheckItemId).
		To(i.reportItem).
		Doc("report an inappropriate item").
		Operation("reportItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")))

	container.Add(ws)
}

//--------------------------------------------------------------------//
// Request Functions

func (i *ItemResource) getFeed(request *restful.Request, response *restful.Response) {
	long, _ := strconv.ParseFloat(request.QueryParameter("longitude"), 64)
	lat, _ := strconv.ParseFloat(request.QueryParameter("latitude"), 64)
	num, _ := strconv.ParseInt(request.QueryParameter("number"), 10, 64)

	res := i.storage.GetFeed(long, lat, int(num))
	var procRes []ItemListPresenter

	for _, r := range *res {
		var im Image
		ims := *i.imageStorage.GetImagesByItemId(r.Id)
		if len(ims) != 0 {
			im = ims[0]
		} else {
			im = Image{}
		}
		procRes = append(procRes, r.ToItemListPresenter(im.ToPresenter(THUMB)))
	}
	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(procRes)

}
func (i *ItemResource) findItem(request *restful.Request, response *restful.Response) {
	id := bson.ObjectIdHex(request.PathParameter("item_id"))
	item := i.storage.GetItem(id)
	ims := *i.imageStorage.GetImagesByItemId(id)
	var ids []bson.ObjectId
	for _, im := range ims {
		ids = append(ids, im.Id)
	}

	pres := item.ToPresenter()
	pres.Images = ids
	response.WriteEntity(pres)
}

func (i *ItemResource) getPicture(request *restful.Request, response *restful.Response) {
	imId := bson.ObjectIdHex(request.PathParameter("image_id"))
	im := i.imageStorage.GetImage(imId)

	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(im.ToPresenter(FULL))

}

func (i *ItemResource) postOffer(request *restful.Request, response *restful.Response) {
	oPres := new(OfferPresenter)
	request.ReadEntity(oPres)
	offer := oPres.ToOffer()
	id := bson.ObjectIdHex(request.PathParameter("item_id"))
	uid := bson.ObjectIdHex(request.PathParameter("user_id"))

	// enforce that we already have one offer existing

	offer.User_Id = uid
	offer.Item_Id = id

	item := i.storage.GetItem(id)
	for _, o := range item.Offers {
		if o.User_Id == uid {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusBadRequest, "You've already made an offer.")
			return
		}
	}
	item.Offers = append(item.Offers, offer)
	i.storage.UpdateItem(*item)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(offer.ToPresenter())
}

func (i *ItemResource) deleteOffer(request *restful.Request, response *restful.Response) {
	id := bson.ObjectIdHex(request.PathParameter("item_id"))
	uid := bson.ObjectIdHex(request.PathParameter("user_id"))

	item := i.storage.GetItem(id)
	var updatedOffers []Offer
	var deletedOffer Offer
	for _, o := range item.Offers {
		if o.User_Id != uid {
			updatedOffers = append(updatedOffers, o)
		} else {
			deletedOffer = o
		}
	}
	item.Offers = updatedOffers
	i.storage.UpdateItem(*item)
	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(deletedOffer.ToPresenter())
}

func (i *ItemResource) reportItem(request *restful.Request, response *restful.Response) {
	id := bson.ObjectIdHex(request.PathParameter("item_id"))

	// eventually have some reporting, maybe by email
	util.Log(id.String() + " was reported as inappropriate")
	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(id.String() + " was reported as inappropriate")
}
