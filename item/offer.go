package item

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/PrincetonOBO/OBOBackend/user"
)

type Offer struct {
	Item_Id       bson.ObjectId      `json:"id" bson:"_id,omitempty"`
	User_Id       bson.ObjectId      `json:"user_id"`
	Price         float64            `json:"price"`
	UserPresenter user.UserPresenter `json:"user"`
}

type OfferPresenter struct {
	Price         float64            `json:"price"`
	UserPresenter user.UserPresenter `json:"user"`
}

func (i Offer) ToPresenter() OfferPresenter {
	return OfferPresenter{Price: i.Price,
		UserPresenter: i.UserPresenter}
}

func (i OfferPresenter) ToOffer() Offer {
	return Offer{Price: i.Price, UserPresenter: i.UserPresenter}
}

func Present(offers []Offer) []OfferPresenter {
	var result []OfferPresenter
	for _, p := range offers {
		result = append(result, p.ToPresenter())
	}
	return result
}
