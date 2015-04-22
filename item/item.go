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
	Longitude   float64       `json:"longitude"`
	Latitude    float64       `json:"latitude"`
}

func (i Item) ToPresenter() ItemPresenter {
	return ItemPresenter{Id: i.Id, Description: i.Description,
		Price: i.Price, Longitude: i.Longitude, Latitude: i.Latitude}
}

func (i *ItemPresenter) ToItem() Item {
	return Item{Id: i.Id, Description: i.Description,
		Price: i.Price, Longitude: i.Longitude, Latitude: i.Latitude}
}

/*
func Present(items []Item) []ItemPresenter {
	var result []ItemPresenter
	for i, p := range items {
		append(result, p.ToPresenter())
	}
	return result
}*/
