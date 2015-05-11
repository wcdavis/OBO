package item

import (
	"github.com/PrincetonOBO/OBOBackend/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// eventually plug into mongodb
type ItemStorage struct {
	db    *mgo.Database   // straight from mongo
	col   *mgo.Collection // collection right from mongo
	scope float64
}

func NewItemStorage(db *mgo.Database) *ItemStorage {
	is := new(ItemStorage)
	is.db = db
	is.col = db.C("item")
	is.scope = 2000.0

	locIndex := mgo.Index{Key: []string{"$2dsphere:location"}}
	textIndex := mgo.Index{Key: []string{"$text:title"}}

	util.Logerr(is.col.EnsureIndex(locIndex))
	util.Logerr(is.col.EnsureIndex(textIndex))
	//util.Logerr(is.col.EnsureIndexKey("title"))
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

func (is *ItemStorage) GetFeed(long float64, lat float64, num int) *[]Item {
	result := []Item{}
	err := is.col.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": is.scope,
			},
		},
	}).Sort("-_id").Limit(num).All(&result)
	util.Logerr(err)
	return &result
}

func (is *ItemStorage) SearchItems(search string, long float64, lat float64, num int) *[]Item {
	result := []Item{}
	err := is.col.Find(bson.M{
		"location": bson.M{
			"$geoWithin": bson.M{
				"$centerSphere": []interface{}{[]float64{long, lat}, is.scope / 6378100.0},
			},
		},

		"$text": bson.M{"$search": search},
	}).Limit(num).All(&result)
	util.Logerr(err)
	return &result
}

func (is *ItemStorage) GetItemsByOffer(user_id bson.ObjectId) *[]Item {
	result := []Item{}
	util.Logerr(is.col.Find(bson.M{"offers": bson.M{"$elemMatch": bson.M{"user_id": user_id}}}).All(&result))
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
