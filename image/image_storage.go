package image

import (
	"github.com/PrincetonOBO/OBOBackend/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// eventually plug into mongodb
type ImageStorage struct {
	db  *mgo.Database   // straight from mongo
	col *mgo.Collection // collection right from mongo
}

func NewImageStorage(db *mgo.Database) *ImageStorage {
	is := new(ImageStorage)
	is.db = db
	is.col = db.C("image")
	return is
}

func (is *ImageStorage) ExistsImage(id bson.ObjectId) bool {
	n, err := is.col.FindId(id).Count()
	util.Logerr(err)
	return n > 0
}

func (is *ImageStorage) GetImage(id bson.ObjectId) *Image {
	result := Image{}
	util.Logerr(is.col.FindId(id).One(&result))
	return &result
}

func (is *ImageStorage) GetImagesByItemId(item_id bson.ObjectId) *[]Image {
	result := []Image{}
	util.Logerr(is.col.Find(bson.M{"item_id": item_id}).All(&result))
	return &result
}

func (is *ImageStorage) InsertImage(image Image) (bool, bson.ObjectId) {
	image.Id = bson.NewObjectId()
	util.Logerr(is.col.Insert(image))
	return true, image.Id
}

func (is *ImageStorage) UpdateImage(image Image) bool {
	util.Logerr(is.col.UpdateId(image.Id, image))
	return true
}

func (is *ImageStorage) DeleteImage(id bson.ObjectId) *Image {
	result := Image{}
	util.Logerr(is.col.FindId(id).One(&result))
	util.Logerr(is.col.RemoveId(id))
	return &result
}

func (is *ImageStorage) Length() int {
	n, err := is.col.Count()
	util.Logerr(err)
	return n
}
