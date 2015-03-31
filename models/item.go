package main

type Item struct {
	id          int
	userId      int
	description string
	price       float64
	longitude   float64
	latitude    float64
}

// implements Cacheable
func (i Item) GetId() int {
	return i.id
}
