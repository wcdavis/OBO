package main

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Description  string        `json:"description"`
	Size         string        `json:"size"`
	Price        string        `json:"price"`
	Latitude     string        `json:"latitude"`
	Longitude    string        `json:"longitude"`
}

// implements Cacheable
func (i Item) GetId() bson.ObjectId {
	return i.Id
}
