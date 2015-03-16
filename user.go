package main

type User struct {
	Id        int
	FirstName string
	LastName  string
	NetId     string
}

// implements Cacheable
func (u User) GetId() int {
	return u.Id
}
