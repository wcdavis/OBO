package item

import (
	"gopkg.in/mgo.v2/bson"

	. "github.com/PrincetonOBO/OBOBackend/image"
)

type Item struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User_Id     bson.ObjectId `json:"user_id"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Offers      []Offer       `json:"offers"`
	Location    GeoJson       `bson:"location" json:"location"`
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type ItemPresenter struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Offers      []Offer       `json:"offers"`
	Location    GeoJson       `bson:"location" json:"location"`
}

type ItemListPresenter struct {
	Id          bson.ObjectId  `json:"id" bson:"_id,omitempty"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Thumbnail   ImagePresenter `json:"thumbnail"`
	Location    GeoJson        `bson:"location" json:"location"`
}

func (i Item) ToPresenter() ItemPresenter {
	return ItemPresenter{Id: i.Id, Description: i.Description,
		Price: i.Price, Offers: nil, Location: i.Location}
}

func (i Item) ToItemListPresenter(im ImagePresenter) ItemListPresenter {
	return ItemListPresenter{Id: i.Id, Description: i.Description,
		Price: i.Price, Thumbnail: im, Location: i.Location}
}

func (i Item) ToPresenterWithOffer(userId bson.ObjectId) ItemPresenter {
	item := ItemPresenter{Id: i.Id, Description: i.Description,
		Price: i.Price, Offers: nil, Location: i.Location}
	for _, o := range i.Offers {
		if o.User_Id == userId {
			item.Offers = append(item.Offers, o)
		}
	}
	return item
}

func (i *ItemPresenter) ToItem() Item {
	return Item{Id: i.Id, Description: i.Description,
		Price: i.Price, Location: i.Location}
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
