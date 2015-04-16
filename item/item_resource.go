package item

import (
	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"
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

	ws.Route(ws.GET("/{item_id}").To(i.findItem).
		Doc("Find an item").
		Operation("findItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Writes(Item{})) // on the response

	ws.Route(ws.PUT("/{item_id}").To(i.updateItem).
		Doc("update an item").
		Operation("updateItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		ReturnsError(409, "duplicate itemId", nil).
		Reads(Item{})) // from the request

	ws.Route(ws.POST("").To(i.createItem).
		Doc("create an item").
		Operation("createItem").
		Reads(Item{})) // from the request

	ws.Route(ws.DELETE("/{item_id}").To(i.removeItem).
		Doc("delete an item").
		Operation("deleteItem").
		Param(ws.PathParameter("item_id", "identifier of the item").DataType("string")).
		Writes(Item{}))

	container.Add(ws)
}

//--------------------------------------------------------------------//
// Request Functions

func (i *ItemResource) findItem(request *restful.Request, response *restful.Response) {
	id, success := i.checkItemId(request, response)
	if !success {
		return
	}
	item := i.storage.GetItem(id)
	response.WriteEntity(item)
}

func (i *ItemResource) createItem(request *restful.Request, response *restful.Response) {
	item, success := i.checkItem(request, response)
	if !success {
		return
	}

	_, item.Id = i.storage.InsertItem(item)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(item)
}

func (i *ItemResource) updateItem(request *restful.Request, response *restful.Response) {
	id, success1 := i.checkItemId(request, response)
	item, success2 := i.checkItem(request, response)
	if !success1 || !success2 {
		return
	}

	item.Id = id // make sure the id is consistent

	i.storage.UpdateItem(item)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(item)
}

func (i *ItemResource) removeItem(request *restful.Request, response *restful.Response) {
	id, success := i.checkItemId(request, response)
	if !success {
		return
	}

	item := i.storage.DeleteItem(id)
	response.WriteHeader(http.StatusAccepted)
	response.WriteEntity(item)
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

func (i *ItemResource) checkItem(request *restful.Request, response *restful.Response) (Item, bool) {
	success := true

	item := new(Item)
	err := request.ReadEntity(item)

	if err != nil {
		success = false
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed item.")
	}

	return *item, success
}
