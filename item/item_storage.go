package item

import (
	"github.com/PrincetonOBO/OBOBackend/util"
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
	util.Logerr(err)
	return n > 0
}

func (is *ItemStorage) GetItem(id bson.ObjectId) *Item {
	result := Item{}
	util.Logerr(is.col.FindId(id).One(&result))
	return &result
}

func (is *ItemStorage) GetItemsByUserId(user_id bson.ObjectId) *[]Item {
	result := []Item{}
	util.Logerr(is.col.Find(bson.M{"user_id": user_id}).All(&result))
	return &result
}

func (is *ItemStorage) InsertItem(item Item) (bool, bson.ObjectId) {
	item.Id = bson.NewObjectId()
	util.Logerr(is.col.Insert(item))
	return true, item.Id
}

func (is *ItemStorage) UpdateItem(item Item) bool {
	util.Logerr(is.col.UpdateId(item.Id, item))
	return true
}

func (is *ItemStorage) DeleteItem(id bson.ObjectId) *Item {
	result := Item{}
	util.Logerr(is.col.FindId(id).One(&result))
	util.Logerr(is.col.RemoveId(id))
	return &result
}

func (is *ItemStorage) Length() int {
	n, err := is.col.Count()
	util.Logerr(err)
	return n
}
