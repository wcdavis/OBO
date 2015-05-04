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

	"strconv"
)

type ItemResource struct {
	storage      *ItemStorage
	userStorage  *UserStorage
	imageStorage *ImageStorage
}

func NewItemResource(db *mgo.Database) *ItemResource {
	ir := new(ItemResource)
	ir.storage = NewItemStorage(db)
	ir.userStorage = NewUserStorage(db)
	ir.imageStorage = NewImageStorage(db)
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

	ws.Route(ws.GET("/").To(i.getFeed).
		Doc("get a feed of items").
		Operation("getFeed").
		Param(ws.QueryParameter("longitude", "longitude for query").DataType("float64")).
		Param(ws.QueryParameter("latitude", "longitude for query").DataType("float64")).
		Param(ws.QueryParameter("number", "number of entries to return").DataType("int")).
		Writes([]ItemListPresenter{}))

	ws.Route(ws.GET("/{item_id}").To(i.findItem).
		Doc("find an item").
		Operation("findItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Writes(ItemPresenter{})) // on the response

	ws.Route(ws.GET("/{item_id}/pic/{pic_id}").To(i.getPicture).
		Doc("get a picture's item").
		Operation("getPic").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Param(ws.PathParameter("pic_id", "identifier of the picture").DataType("string")).
		Writes(ImagePresenter{})) // from the request

	ws.Route(ws.POST("/{item_id}/offer/users/{user_id}").To(i.postOffer).
		Doc("post an offer").
		Operation("newOffer").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		// may be replaced by user token
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Reads(OfferPresenter{})) // from the request

	ws.Route(ws.DELETE("/{item_id}/offer/users/{user_id}").To(i.deleteOffer).
		Doc("delete an offer").
		Operation("deleteOffer").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		// may be replaced
		Param(ws.PathParameter("user_id", "identifier of the user").DataType("string")).
		Writes(OfferPresenter{})) // from the request

	ws.Route(ws.GET("/{item_id}/report").To(i.reportItem).
		Doc("report an inappropriate item").
		Operation("reportItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")))

	container.Add(ws)
}

//--------------------------------------------------------------------//
// Request Functions

func (i *ItemResource) getFeed(request *restful.Request, response *restful.Response) {
	// this is where we would do the geo query, but right now we don't
	long, e1 := strconv.ParseFloat(request.QueryParameter("longitude"), 64)
	lat, e2 := strconv.ParseFloat(request.QueryParameter("latitude"), 64)
	num, e3 := strconv.ParseInt(request.QueryParameter("number"), 10, 64)

	if e1 != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed longitude.")
		return
	}
	if e2 != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed latitude.")
		return
	}
	if e3 != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed number.")
		return
	}

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
	id, success := i.checkItemId(request, response)
	if !success {
		return
	}
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
	imId, success1 := i.checkImageId(request, response)
	itemId, success2 := i.checkItemId(request, response)

	if !success1 || !success2 {
		return
	}

	im := i.imageStorage.GetImage(imId)

	if im.Item_Id != itemId {
		return
	}
	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(im.ToPresenter(FULL))

}

func (i *ItemResource) postOffer(request *restful.Request, response *restful.Response) {
	offer, success1 := i.checkOffer(request, response)
	id, success2 := i.checkItemId(request, response)
	uid, success3 := i.checkUserId(request, response)
	if !success1 || !success2 || !success3 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Incorrect image id.")
	}

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
	id, success2 := i.checkItemId(request, response)
	uid, success3 := i.checkUserId(request, response)
	if !success2 || !success3 {
		return
	}

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
	id, success2 := i.checkItemId(request, response)
	if !success2 {
		return
	}
	util.Log(id.String() + " was reported as inappropriate")
	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(id.String() + " was reported as inappropriate")
}

//--------------------------------------------------------------------//
// Utility Functions

func (i *ItemResource) checkItemId(request *restful.Request, response *restful.Response) (bson.ObjectId, bool) {
	success := true
	idString := request.PathParameter("item_id")

	if !bson.IsObjectIdHex(idString) {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed item id.")
	} else if !i.storage.ExistsItem(bson.ObjectIdHex(idString)) {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Item not found.")
	}
	id := bson.ObjectIdHex(idString)

	return id, success
}

func (i *ItemResource) checkOffer(request *restful.Request, response *restful.Response) (Offer, bool) {
	success := true

	offer := new(OfferPresenter)
	err := request.ReadEntity(offer)
	util.Logerr(err)

	if err != nil {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed offer.")
	}

	return (offer.ToOffer()), success
}

func (i *ItemResource) checkImageId(request *restful.Request, response *restful.Response) (bson.ObjectId, bool) {
	success := true
	idString := request.PathParameter("pic_id")

	if !bson.IsObjectIdHex(idString) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed pic_id.")
		return bson.NewObjectId(), false
	} else if !i.imageStorage.ExistsImage(bson.ObjectIdHex(idString)) {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Image not found.")
	}
	id := bson.ObjectIdHex(idString)

	return id, success
}

func (i *ItemResource) checkUserId(request *restful.Request, response *restful.Response) (bson.ObjectId, bool) {
	success := true
	idString := request.PathParameter("user_id")

	if !bson.IsObjectIdHex(idString) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed user_id.")
		return bson.NewObjectId(), false
	} else if !i.userStorage.ExistsUser(bson.ObjectIdHex(idString)) {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "User not found.")
	}
	id := bson.ObjectIdHex(idString)

	return id, success
}
