package main

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	NetId     string        `json:"net_id"`
}

// implements Cacheable
func (u User) GetId() bson.ObjectId {
	return u.Id
}
