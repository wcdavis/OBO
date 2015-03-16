package main

// we want to implement a generic-ish cache
type Cacheable interface {
	GetId() int
}

type Cache struct {
	m map[int]Cacheable
}

func NewCache() *Cache {
	c := new(Cache)
	c.m = make(map[int]Cacheable)
	return c
}

func (c *Cache) contains(id int) bool {
	_, ok := c.m[id]
	return ok
}

func (c *Cache) get(id int) Cacheable {
	obj, _ := c.m[id]
	return obj
}

func (c *Cache) insert(obj Cacheable) bool {
	c.m[obj.GetId()] = obj
	return true
}

func (c *Cache) remove(id int) Cacheable {
	obj, _ := c.m[id]
	delete(c.m, id)
	return obj
}

func (c *Cache) length() int {
	return len(c.m)
}
