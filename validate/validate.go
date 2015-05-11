package validate

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"strconv"

	. "github.com/PrincetonOBO/OBOBackend/image"
	. "github.com/PrincetonOBO/OBOBackend/item"
	. "github.com/PrincetonOBO/OBOBackend/user"

	"github.com/PrincetonOBO/OBOBackend/util"
)

var (
	privateKey []byte
	publicKey  []byte
)

func init() {
	publicKey, _ = ioutil.ReadFile("public.pem")
	privateKey, _ = ioutil.ReadFile("private.new.pem")
}

type Validator struct {
	itemStorage  *ItemStorage
	userStorage  *UserStorage
	imageStorage *ImageStorage
}

func NewValidator(db *mgo.Database) *Validator {
	val := new(Validator)
	val.itemStorage = NewItemStorage(db)
	val.userStorage = NewUserStorage(db)
	val.imageStorage = NewImageStorage(db)
	return val
}

func (v *Validator) CheckItem(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	itemPres := new(ItemPresenter)
	err := request.ReadEntity(itemPres)
	util.Logerr(err)

	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed item.")
		return
	}
	if !checkCoordinates(itemPres.Location.Coordinates) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed coordinates. Should be [longitude,latitude].")
		return
	}
	chain.ProcessFilter(request, response)
}

func checkCoordinates(coords []float64) bool {
	if len(coords) != 2 {
		return false
	}
	if coords[0] > 180 || coords[0] < -180 {
		return false
	}
	if coords[1] > 90 || coords[1] < -90 {
		return false
	}
	return true
}

func (v *Validator) CheckItemId(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	idString := request.PathParameter("item_id")

	if !bson.IsObjectIdHex(idString) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed item id.")
		return
	} else if !v.itemStorage.ExistsItem(bson.ObjectIdHex(idString)) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Item not found.")
		return
	}
	chain.ProcessFilter(request, response)
}

func (v *Validator) CheckImage(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	imagePres := new(ImagePresenter)
	err1 := request.ReadEntity(imagePres)

	if err1 != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed image.")
		return
	}

	var imBytes []byte
	imBytes, err2 := base64.StdEncoding.DecodeString(imagePres.Image)
	if err2 != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed base64 encoding.")
		return
	}

	buf := bytes.NewBuffer(imBytes)
	_, err3 := jpeg.Decode(buf)
	if err3 != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed jpeg image.")
		return
	}
	chain.ProcessFilter(request, response)
}

func (v *Validator) CheckImageId(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	idString := request.PathParameter("pic_id")

	if !bson.IsObjectIdHex(idString) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed pic_id.")
		return
	} else if !v.imageStorage.ExistsImage(bson.ObjectIdHex(idString)) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Image not found.")
		return
	}
	chain.ProcessFilter(request, response)
}

func (v *Validator) CheckUser(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	usr := new(User)
	err := request.ReadEntity(usr)

	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed user.")
		return
	}
	chain.ProcessFilter(request, response)
}

func (v *Validator) CheckUserId(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	idString := request.PathParameter("user_id")
	if !bson.IsObjectIdHex(idString) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed user_id.")
		return
	} else if !v.userStorage.ExistsUser(bson.ObjectIdHex(idString)) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "User not found.")
		return
	}
	chain.ProcessFilter(request, response)

}

func (v *Validator) CheckFeedQuery(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	_, e1 := strconv.ParseFloat(request.QueryParameter("longitude"), 64)
	_, e2 := strconv.ParseFloat(request.QueryParameter("latitude"), 64)
	_, e3 := strconv.ParseInt(request.QueryParameter("number"), 10, 64)

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
	chain.ProcessFilter(request, response)
}

func (v *Validator) CheckOffer(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	offer := new(OfferPresenter)
	err := request.ReadEntity(offer)

	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed offer.")
		return
	}
	chain.ProcessFilter(request, response)
}

func (v *Validator) CheckOfferId(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	idString := request.PathParameter("offer_id")

	if !bson.IsObjectIdHex(idString) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Malformed offer id.")
		return
	}
	chain.ProcessFilter(request, response)
}

func (v *Validator) CheckItemOwnership(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	id := bson.ObjectIdHex(request.PathParameter("item_id"))
	uid := bson.ObjectIdHex(request.PathParameter("user_id"))

	item := v.itemStorage.GetItem(id)
	if item.User_Id != uid {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "User doesn't own item")
		return
	}
	chain.ProcessFilter(request, response)
}

func (v *Validator) Authenticate(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	encoded := request.Request.Header.Get("Authorization")
	token, err := jwt.Parse(encoded, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	validId := true

	if request.PathParameter("user_id") != "" && token != nil {
		uid := request.PathParameter("user_id")
		if uid != token.Claims["user_id"] {
			validId = false
		}
	}
	if err == nil && token.Valid && validId {
		chain.ProcessFilter(request, response)
	} else {
		response.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		response.WriteErrorString(401, "401: Not Authorized")
		return
	}

}

func (v *Validator) CreateAuthenticatedToken(user User) string {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims["user_id"] = user.Id.Hex()
	tokenString, _ := token.SignedString(privateKey)

	return tokenString
}
