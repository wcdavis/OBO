package main

import (
	"gopkg.in/mgo.v2"
)

// eventually plug into mongodb
type UserStorage struct {
	db    *mgo.Database   // straight from mongo
	col   *mgo.Collection // collection right from mongo
	cache *Cache          // right now we map a uid to User
}

func NewUserStorage(db *mgo.Database) *UserStorage {
	us := new(UserStorage)
	us.cache = NewCache()
	us.col = db.C("user")
	return us
}

func (u *UserStorage) ExistsUser(id int) bool {
	return u.cache.contains(id)
}

func (u *UserStorage) GetUser(id int) User {
	user, _ := u.cache.get(id).(User)
	return user

}

func (u *UserStorage) InsertUser(user User) bool {
	u.cache.insert(user)
	return true
}

func (u *UserStorage) UpdateUser(user User) User {
	u.cache.insert(user)
	u1, _ := u.cache.get(user.GetId()).(User)
	return u1
}

func (u *UserStorage) DeleteUser(id int) User {
	user, _ := u.cache.remove(id).(User)
	return user
}

func (u *UserStorage) Length() int {
	return u.cache.length()
}
