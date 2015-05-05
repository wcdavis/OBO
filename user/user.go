package user

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name           string        `json:"name"`
	NetId          string        `json:"net_id"`
	PickupLocation string        `json:"pickup_loc"`
}

type UserPresenter struct {
	Name           string `json:"name"`
	NetId          string `json:"net_id"`
	PickupLocation string `json:"pickup_loc"`
}

func (u User) ToPresenter() UserPresenter {
	return UserPresenter{Name: u.Name,
		NetId: u.NetId, PickupLocation: u.PickupLocation}
}

func (u *UserPresenter) ToUser() User {
	return User{Name: u.Name, NetId: u.NetId, PickupLocation: u.PickupLocation}
}
