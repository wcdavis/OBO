package item

import (
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"

	"github.com/PrincetonOBO/OBOBackend/util"

	"strconv"
)

type ItemResource struct {
	storage *ItemStorage
}

func NewItemResource(db *mgo.Database) *ItemResource {
	ir := new(ItemResource)
	ir.storage = NewItemStorage(db)
	return ir
}

// significant boilerplate for registration adapted from
// https://github.com/emicklei/go-restful/blob/master/examples/restful-user-resource.go
func (i ItemResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/items").
		Doc("Manage Items").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/").To(i.getFeed).
		Doc("Get a feed of items").
		Operation("getFeed").
		Param(ws.QueryParameter("longitude", "longitude for query").DataType("float64")).
		Param(ws.QueryParameter("latitude", "longitude for query").DataType("float64")).
		Param(ws.QueryParameter("number", "number of entries to return").DataType("int")).
		Writes([]ItemPresenter{}))

	ws.Route(ws.GET("/{item_id}").To(i.findItem).
		Doc("Find an item").
		Operation("findItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Writes(ItemPresenter{})) // on the response

	ws.Route(ws.POST("/{item_id}/offer").To(i.postOffer).
		Doc("post an offer").
		Operation("newOffer").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		ReturnsError(409, "duplicate itemId", nil).
		Reads(OfferPresenter{})) // from the request

	ws.Route(ws.DELETE("/{item_id}/offer").To(i.deleteOffer).
		Doc("delete an offer").
		Operation("deleteOffer").
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
	num, e3 := strconv.ParseFloat(request.QueryParameter("number"), 64)

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

	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity([5]ItemPresenter{Id: bson.ObjectIdHex("5536a8dd66580e3d7e000001"),
		Description: "here's another item",
		Price:       55.1,
		Longitude:   23.4,
		Latitude:    42.3})

}
func (i *ItemResource) findItem(request *restful.Request, response *restful.Response) {
	id, success := i.checkItemId(request, response)
	if !success {
		return
	}
	item := i.storage.GetItem(id)
	response.WriteEntity(item.ToPresenter())
}

func (i *ItemResource) postOffer(request *restful.Request, response *restful.Response) {
	offer, success1 := i.checkOffer(request, response)
	id, success2 := i.checkItemId(request, response)
	if !success1 || !success2 {
		return
	}

	// enforce that we already have one offer existing

	item := i.storage.GetItem(id)
	for i, o := range item.Offers {
		if o.Item_Id == id {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusBadRequest, "You've already made an offer.")
			return
		}
	}
	append(item.Offers, offer)
	i.storage.UpdateItem(item)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(offer.ToPresenter())
}

func (i *ItemResource) deleteOffer(request *restful.Request, response *restful.Response) {
	id, success2 := i.checkItemId(request, response)
	if !success2 {
		return
	}

	item := i.storage.GetItem(id)
	var updatedOffers []Offer
	var deletedOffer Offer
	for i, o := range item.Offers {
		if o.Item_Id != id {
			append(updatedOffers, o)
		} else {
			deletedOffer = o
		}
	}
	item.Offers = updatedOffers
	i.storage.UpdateItem(item)
	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(offer.ToPresenter())
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

	return *(offer.ToOffer()), success
}
