package item

import (
	"gopkg.in/mgo.v2/bson"
)

type Item struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User_Id     bson.ObjectId `json:"user_id"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Offers      []Offer       `json:"offers"`
	Longitude   float64       `json:"longitude"`
	Latitude    float64       `json:"latitude"`
}

type ItemPresenter struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Offers      []Offer       `json:"offers"`
	Longitude   float64       `json:"longitude"`
	Latitude    float64       `json:"latitude"`
}

func (i Item) ToPresenter() ItemPresenter {
	return ItemPresenter{Id: i.Id, Description: i.Description,
		Price: i.Price, Offers: nil, Longitude: i.Longitude, Latitude: i.Latitude}
}

func (i Item) ToPresenterWithOffer(userId bson.ObjectId) ItemPresenter {
	item := ItemPresenter{Id: i.Id, Description: i.Description,
		Price: i.Price, Offers: nil, Longitude: i.Longitude, Latitude: i.Latitude}
	for _, o := range i.Offers {
		if o.User_Id == userId {
			item.Offers = append(item.Offers, o)
		}
	}
	return item
}

func (i *ItemPresenter) ToItem() Item {
	return Item{Id: i.Id, Description: i.Description,
		Price: i.Price, Longitude: i.Longitude, Latitude: i.Latitude}
}

func Present(items []Item) []ItemPresenter {
	var result []ItemPresenter
	for _, p := range items {
		result = append(result, p.ToPresenter())
	}
	return result
}

func PresentWithOffer(items []Item, userId bson.ObjectId) []ItemPresenter {
	var result []ItemPresenter
	for _, p := range items {
		result = append(result, p.ToPresenterWithOffer(userId))
	}
	return result
}
