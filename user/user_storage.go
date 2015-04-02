package user

import (
	"github.com/PrincetonOBO/OBOBackend/util"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// eventually plug into mongodb
type UserStorage struct {
	db  *mgo.Database   // straight from mongo
	col *mgo.Collection // collection right from mongo
}

func newUserStorage(db *mgo.Database) *UserStorage {
	us := new(UserStorage)
	us.db = db
	us.col = db.C("user")
	return us
}

func (u *UserStorage) ExistsUser(id bson.ObjectId) bool {
	n, err := u.col.FindId(id).Count()
	util.Logerr(err)
	return n > 0
}

func (u *UserStorage) GetUser(id bson.ObjectId) *User {
	result := User{}
	util.Logerr(u.col.FindId(id).One(&result))
	return &result
}

func (u *UserStorage) InsertUser(user User) (bool, bson.ObjectId) {
	user.Id = bson.NewObjectId()
	util.Logerr(u.col.Insert(user))
	return true, user.Id
}

func (u *UserStorage) UpdateUser(user User) bool {
	util.Logerr(u.col.UpdateId(user.Id, user))
	return true
}

func (u *UserStorage) DeleteUser(id bson.ObjectId) *User {
	result := User{}
	util.Logerr(u.col.FindId(id).One(&result))
	util.Logerr(u.col.RemoveId(id))
	return &result
}

func (u *UserStorage) Length() int {
	n, err := u.col.Count()
	util.Logerr(err)
	return n
}
