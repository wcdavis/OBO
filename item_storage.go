package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// eventually plug into mongodb
type ItemStorage struct {
	db  *mgo.Database   // straight from mongo
	col *mgo.Collection // collection right from mongo
}

func NewItemStorage(db *mgo.Database) *ItemStorage {
	is := new(ItemStorage)
	is.db = db
	is.col = db.C("item")
	return is
}

func (is *ItemStorage) ExistsItem(id bson.ObjectId) bool {
	n, err := is.col.FindId(id).Count()
	logerr(err)
	return n > 0
}

func (is *ItemStorage) GetItem(id bson.ObjectId) *User {
	result := Item{}
	logerr(is.col.FindId(id).One(&result))
	return &result
}

func (is *ItemStorage) InsertItem(item Item) (bool, bson.ObjectId) {
	item.Id = bson.NewObjectId()
	logerr(is.col.Insert(user))
	return true, item.Id
}

func (is *ItemStorage) UpdateItem(item Item) bool {
	logerr(is.col.UpdateId(item.Id, user))
	return true
}

func (is *UserStorage) DeleteUser(id bson.ObjectId) *User {
	result := Item{}
	logerr(is.col.FindId(id).One(&result))
	logerr(is.col.RemoveId(id))
	return &result
}

func (is *UserStorage) Length() int {
	n, err := is.col.Count()
	logerr(err)
	return n
}
