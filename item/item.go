package item

import (
	"gopkg.in/mgo.v2/bson"
)

type Item struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Description string        `json:"description"`
	price       float64       `json:"price"`
	longitude   float64       `json:"longitude"`
	latitude    float64       `json:"latitude"`
}

// implements Cacheable
func (i Item) GetId() bson.ObjectId {
	return i.Id
}
