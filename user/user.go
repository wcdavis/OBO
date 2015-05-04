package user

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	FirstName      string        `json:"first_name"`
	LastName       string        `json:"last_name"`
	NetId          string        `json:"net_id"`
	PickupLocation string        `json:"pickup_loc"`
}

type UserPresenter struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	NetId          string `json:"net_id"`
	PickupLocation string `json:"pickup_loc"`
}

func (u User) ToPresenter() UserPresenter {
	return UserPresenter{FirstName: u.FirstName, LastName: u.LastName,
		NetId: u.NetId, PickupLocation: u.PickupLocation}
}

func (u *UserPresenter) ToUser() User {
	return User{FirstName: u.FirstName, LastName: u.LastName,
		NetId: u.NetId, PickupLocation: u.PickupLocation}
}
