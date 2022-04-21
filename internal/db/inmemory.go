package db

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	HashCollection      = "hashCollection"
	ResourcesCollection = "resourcesCollection"
)

var (
	colExistErr    = errors.New("collection already exist")
	colNotExistErr = errors.New("collection not exist")
	keyExistErr    = errors.New("key already exist")
	keyNotExistErr = errors.New("key not exist")
)

type InMemoryDB struct {
	cols map[string]*InMemoryCollection
	m    *sync.Mutex
}

type InMemoryCollection struct {
	vals map[string]interface{}
	m    *sync.Mutex
}

func InitInMemoryDB() (*InMemoryDB, error) {
	db := &InMemoryDB{
		m:    &sync.Mutex{},
		cols: make(map[string]*InMemoryCollection),
	}
	_, err := db.CreateCollection(HashCollection)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"collection": HashCollection,
		}).Panic("cannot create collection")
		return nil, err
	}
	_, err = db.CreateCollection(ResourcesCollection)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"collection": ResourcesCollection,
		}).Panic("cannot create collection")
		return nil, err
	}
	return db, err
}

func (i *InMemoryDB) CreateCollection(n string) (Collection, error) {
	i.m.Lock()
	defer i.m.Unlock()
	if i.CollectionExist(n) {
		return nil, colExistErr
	}
	i.cols[n] = &InMemoryCollection{
		vals: make(map[string]interface{}),
		m:    &sync.Mutex{},
	}

	return i.cols[n], nil
}

func (i *InMemoryDB) GetCollection(n string) (Collection, error) {
	i.m.Lock()
	defer i.m.Unlock()
	if !i.CollectionExist(n) {
		return nil, colNotExistErr
	}
	return i.cols[n], nil
}

func (i *InMemoryDB) DeleteCollection(n string) error {
	i.m.Lock()
	defer i.m.Unlock()
	if !i.CollectionExist(n) {
		return colNotExistErr
	}
	delete(i.cols, n)
	return nil
}

func (i *InMemoryDB) CollectionExist(n string) bool {
	if _, ok := i.cols[n]; ok {
		return true
	}
	return false
}

func (c *InMemoryCollection) AddKeyValue(k string, v interface{}) error {
	c.m.Lock()
	defer c.m.Unlock()
	if c.KeyExist(k) {
		return keyExistErr
	}
	c.vals[k] = v
	return nil
}

func (c *InMemoryCollection) DeleteKey(n string) error {
	c.m.Lock()
	defer c.m.Unlock()
	if !c.KeyExist(n) {
		return keyNotExistErr
	}
	delete(c.vals, n)
	return nil
}

func (c *InMemoryCollection) KeyExist(n string) bool {
	if _, ok := c.vals[n]; ok {
		return true
	}
	return false
}
