package migration

import (
	"context"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Migrations = []app.List{}

func Run() {

	// agregue las funciones de migracion up y down
	add("create permissions", PermissionsUp, PermissionsDown)
	add("create roles", RolesUp, RolesDown)
	add("create users", UsersUp, UsersDown)

}

func add(name string, up func(), down func()) {
	migration := app.List{}
	migration.Set("name", name)
	migration.Set("up", up)
	migration.Set("down", down)
	Migrations = append(Migrations, migration)
}

func CreateIndex(collection string, order int, fields ...string) {
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: order})
	}
	name, er := app.DB.Collection(collection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: keys,
	})
	if er != nil {
		app.PrintError("Failed to create index :collection :error ", app.E("collection", collection), app.E("error", er.Error()))
		panic(er.Error())
	}
	app.PrintInfo("Created index :collection :name", app.E("collection", collection), app.E("name", name))
}

func CreateUniqueIndex(collection string, sort int, fields ...string) {
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: sort})
	}
	name, er := app.DB.Collection("users").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	})
	if er != nil {
		app.PrintError("Failed to create unique index :collection :error ", app.E("collection", collection), app.E("error", er.Error()))
		panic(er.Error())
	}
	app.PrintInfo("Created unique index :collection :name", app.E("collection", collection), app.E("name", name))
}

func CreateIndexWithOptions(collection string, keys bson.D, options *options.IndexOptionsBuilder) {
	name, er := app.DB.Collection(collection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    keys,
		Options: options,
	})
	if er != nil {
		app.PrintError("Failed to create index :collection :error ", app.E("collection", collection), app.E("error", er.Error()))
		panic(er.Error())
	}
	app.PrintInfo("Created index :collection :name", app.E("collection", collection), app.E("name", name))
}

func CreateCollection(collection string, opts ...options.Lister[options.CreateCollectionOptions]) {
	er := app.DB.CreateCollection(context.TODO(), collection, opts...)
	if er != nil {
		app.PrintError("Failed to create collection :collection :error ", app.E("collection", collection), app.E("error", er.Error()))
		panic(er.Error())
	}
	app.PrintInfo("Created collection :collection", app.E("collection", collection))
}

func DropIndex(collection string, indexName string) {

	er := app.DB.Collection(collection).Indexes().DropOne(context.TODO(), indexName)
	if er != nil {
		app.PrintError("Failed to drop index :collection :name :error ", app.E("collection", collection), app.E("name", indexName), app.E("error", er.Error()))
		panic(er.Error())
	}
	app.PrintInfo("Dropped index :collection :name", app.E("collection", collection), app.E("name", indexName))
}

func DropAllIndexes(collection string) {

	er := app.DB.Collection(collection).Indexes().DropAll(context.TODO())
	if er != nil {
		app.PrintError("Failed to drop all indexes :collection :error ", app.E("collection", collection), app.E("error", er.Error()))
		panic(er.Error())
	}
	app.PrintInfo("Dropped all indexes :collection", app.E("collection", collection))
}

func DropCollection(collection string) {

	er := app.DB.Collection(collection).Drop(context.TODO())
	if er != nil {
		app.PrintError("Failed to drop collection :collection :error ", app.E("collection", collection), app.E("error", er.Error()))
		panic(er.Error())
	}
	app.PrintInfo("Dropped collection :collection", app.E("collection", collection))
}
