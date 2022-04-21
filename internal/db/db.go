package db

type DB interface {
	CreateCollection(n string) (Collection, error)
	GetCollection(n string) (Collection, error)
	CollectionExist(n string) bool
	DeleteCollection(n string) error
}

type Collection interface {
	KeyExist(k string) bool
	AddKeyValue(k string, v interface{}) error
	DeleteKey(k string) error
}
