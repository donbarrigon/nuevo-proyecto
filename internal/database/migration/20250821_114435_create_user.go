package migration

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func UsersUp() {
	CreateCollection("users")
	CreateUniqueIndex("users", 1, "email")
	CreateUniqueIndex("users", 1, "profile.nickname")
	CreateIndexWithOptions("users", bson.D{{Key: "deleted_at", Value: 1}}, options.Index().SetSparse(true))
}

func UsersDown() {
	DropCollection("users")
}
