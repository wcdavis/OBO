package item

import (
	"gopkg.in/mgo.v2/bson"
)

type Offer struct {
	Item_Id   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User_Id   bson.ObjectId `json:"user_id"`
	Price     float64       `json:"price"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	NetId     string        `json:"net_id"`
}

type OfferPresenter struct {
	Item_Id   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Price     float64       `json:"price"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	NetId     string        `json:"net_id"`
}

func (i Offer) ToPresenter() OfferPresenter {
	return OfferPresenter{ItemId: i.Item_Id, Price: i.Price,
		FirstName: i.FirstName, LastName: i.LastName, NetId: i.NetId}
}

func (i OfferPresenter) ToOffer() Offer {
	return Offer{ItemId: i.Item_Id, Price: i.Price,
		FirstName: i.FirstName, LastName: i.LastName, NetId: i.NetId}
}

func Present(offers []Offer) []OfferPresenter {
	var result []OfferPresenter
	for i, p := range offers {
		append(result, p.ToPresenter())
	}
	return result
}
