package migration

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func UsersUp() {
	CreateCollection("users", func(collection string) {
		CreateUniqueIndex(collection, 1, "email")
		CreateUniqueIndex(collection, 1, "profile.nickname")
		CreateIndexWithOptions(collection, bson.D{{Key: "deleted_at", Value: 1}}, options.Index().SetSparse(true))

	})
}

func UsersDown() {
	DropCollection("users")
}
