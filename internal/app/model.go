package app

import "go.mongodb.org/mongo-driver/v2/bson"

type Model interface {
	CollectionName() string
	GetID() bson.ObjectID
}

type Migration interface {
	Index() map[string]string
}
