package item

import (
	"gopkg.in/mgo.v2/bson"

	. "github.com/PrincetonOBO/OBOBackend/image"
)

type Item struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User_Id     bson.ObjectId `json:"user_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Price       int           `json:"price"`
	Offers      []Offer       `json:"offers"`
	Location    GeoJson       `bson:"location" json:"location"`
	Time        int64         `json:"time"`
	Sold        bool          `json:"sold"`
	Size        string        `json:"size"`
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type ItemPresenter struct {
	Id          bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       int             `json:"price"`
	Offers      []Offer         `json:"offers"`
	Images      []bson.ObjectId `json:"images"`
	Location    GeoJson         `bson:"location" json:"location"`
	Time        int64           `json:"time"`
	Sold        bool            `json:"sold"`
	Size        string          `json:"size"`
}

type ItemListPresenter struct {
	Id          bson.ObjectId  `json:"id" bson:"_id,omitempty"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Price       int            `json:"price"`
	Thumbnail   ImagePresenter `json:"thumbnail"`
	Location    GeoJson        `bson:"location" json:"location"`
	Time        int64          `json:"time"`
	Sold        bool           `json:"sold"`
	Size        string         `json:"size"`
}

func (i Item) ToPresenter() ItemPresenter {
	return ItemPresenter{Id: i.Id, Title: i.Title, Description: i.Description,
		Price: i.Price, Offers: nil, Images: nil, Location: i.Location, Time: i.Time,
		Sold: i.Sold, Size: i.Size}
}

func (i Item) ToItemListPresenter(im ImagePresenter) ItemListPresenter {
	return ItemListPresenter{Id: i.Id, Description: i.Description, Title: i.Title,
		Price: i.Price, Thumbnail: im, Location: i.Location, Time: i.Time, Sold: i.Sold,
		Size: i.Size}
}

func (i Item) ToPresenterWithOffer(userId bson.ObjectId) ItemPresenter {
	item := ItemPresenter{Id: i.Id, Description: i.Description, Title: i.Title,
		Price: i.Price, Offers: nil, Location: i.Location, Time: i.Time, Sold: i.Sold,
		Size: i.Size}
	for _, o := range i.Offers {
		if o.User_Id == userId {
			item.Offers = append(item.Offers, o)
		}
	}
	return item
}

func (i *ItemPresenter) ToItem() Item {
	return Item{Id: i.Id, Description: i.Description, Title: i.Title,
		Price: i.Price, Location: i.Location, Size: i.Size}
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
