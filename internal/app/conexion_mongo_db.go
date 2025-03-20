package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Mongo *ConexionMongoDB

type ConexionMongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (db *ConexionMongoDB) FindByHexID(model Model, id string) error {

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID [" + id + "] invalid")
	}
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	if err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("Not found [" + id + "]")
		}
		return err
	}
	return nil
}

func (db *ConexionMongoDB) FindByID(model Model, id bson.ObjectID) error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	if err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		// if err == mongo.ErrNoDocuments {
		// 	return nil
		// }
		return err
	}
	return nil
}

func (db *ConexionMongoDB) FindByField(model Model, field string, value any) (*[]Model, error) {
	filter := bson.D{bson.E{Key: field, Value: value}}
	cursor, err := db.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	results := toSlice(model)
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return &results, nil
}

func (db *ConexionMongoDB) FindOneByField(model Model, field string, value any) error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	if err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		return err
	}
	return nil
}

func (db *ConexionMongoDB) FindOne(model Model, filter bson.D) error {
	return db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model)
}

func (db *ConexionMongoDB) Find(model Model, filter bson.D) (*[]Model, error) {
	cursor, err := db.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	results := toSlice(model)
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return &results, nil
}

func toSlice[T any](v T) []T {
	return make([]T, 0)
}

func (db *ConexionMongoDB) Create(model Model) (*mongo.InsertOneResult, error) {

	collection := db.Database.Collection(model.CollectionName())
	return collection.InsertOne(context.TODO(), model)
}

func (db *ConexionMongoDB) Update(model Model) (*mongo.UpdateResult, error) {
	collection := db.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: model}}
	return collection.UpdateOne(context.TODO(), filter, update)
}

func (db *ConexionMongoDB) Destroy(model Model) (*mongo.DeleteResult, error) {
	collection := db.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	return collection.DeleteOne(context.TODO(), filter)
}

func InitMongoDB() error {
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMinPoolSize(5)
	clientOptions.SetRetryWrites(true)
	clientOptions.SetTimeout(30 * time.Second)

	var err error
	Mongo = &ConexionMongoDB{}
	Mongo.Client, err = mongo.Connect(clientOptions)
	if err != nil {
		log.Fatalf("Error al conectar con mongodb: %v", err)
		return err
	}
	Mongo.Database = Mongo.Client.Database(DB_NAME)

	log.Printf("Conectado exitosamente a MongoDB: %s - Base de datos: %s", MONGO_URI, DB_NAME)
	return nil
}

func CloseMongoConnection() error {

	if Mongo.Client == nil {
		return nil
	}

	err := Mongo.Client.Disconnect(context.TODO())
	if err != nil {
		return fmt.Errorf("error al cerrar la conexión con MongoDB: %w", err)
	}

	log.Println("Conexión a MongoDB cerrada correctamente")
	return nil
}
