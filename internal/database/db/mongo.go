package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/config"
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Mongo *ConexionMongoDB

type MongoModel interface {
	CollectionName() string
	GetID() bson.ObjectID
	SetID(id bson.ObjectID)
	Default()
	Validate(lang string) errors.Error
}

type ConexionMongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func FindByHexID(model MongoModel, id string) errors.Error {

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.HexID(err)
	}
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	if err := Mongo.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func FindByID(model MongoModel, id bson.ObjectID) errors.Error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	if err := Mongo.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func FindManyByField(model MongoModel, result any, field string, value any) errors.Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	cursor, err := Mongo.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func FindOneByField(model MongoModel, field string, value any) errors.Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	if err := Mongo.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func FindAll(model MongoModel, result any) errors.Error {

	cursor, err := Mongo.Database.Collection(model.CollectionName()).Find(context.TODO(), bson.D{})
	if err != nil {
		return errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return errors.Mongo(err)
	}
	return nil

}

func FindOne(model MongoModel, filter bson.D) errors.Error {
	err := Mongo.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model)
	if err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func FindMany(model MongoModel, result any, filter bson.D) errors.Error {
	cursor, err := Mongo.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func Create(model MongoModel) errors.Error {
	model.Default()
	collection := Mongo.Database.Collection(model.CollectionName())
	result, err := collection.InsertOne(context.TODO(), model)
	if err != nil {
		return errors.Mongo(err)
	}
	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return errors.Unknown(errors.New("No se logro hacer la conversion a bson.ObjectID"))
	}
	model.SetID(id)
	return nil
}

func Update(model MongoModel) errors.Error {
	model.Default()
	collection := Mongo.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: model}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Mongo(err)
	}

	if result.MatchedCount == 0 {
		return errors.NoDocuments(errors.New("mongo.UpdateResult.MatchedCount == 0"))
	}

	if result.ModifiedCount == 0 {
		return errors.Update(errors.New("mongo.UpdateResult.ModifiedCount == 0"))
	}

	return nil
}
func Delete(model MongoModel) errors.Error {
	collection := Mongo.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{{Key: "deletedAt", Value: time.Now()}}}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Mongo(err)
	}

	if result.MatchedCount == 0 {
		return errors.NoDocuments(errors.New("mongo.UpdateResult.MatchedCount == 0"))
	}

	if result.ModifiedCount == 0 {
		return errors.Delete(errors.New("mongo.UpdateResult.ModifiedCount == 0"))
	}
	return nil
}

func Restore(model MongoModel) errors.Error {
	collection := Mongo.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$unset", Value: bson.D{{Key: "deletedAt", Value: nil}}}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Mongo(err)
	}

	if result.MatchedCount == 0 {
		return errors.NoDocuments(errors.New("mongo.UpdateResult.MatchedCount == 0"))
	}

	if result.ModifiedCount == 0 {
		return errors.Restore(errors.New("mongo.UpdateResult.ModifiedCount == 0"))
	}
	return nil
}

func ForceDelete(model MongoModel) errors.Error {
	collection := Mongo.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return errors.Mongo(err)
	}

	if result.DeletedCount == 0 {
		return errors.ForceDelete(errors.New("mongo.DeleteResult.DeletedCount == 0"))
	}
	return nil
}

func InitMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(config.MONGO_URI).SetServerAPIOptions(serverAPI)
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
	Mongo.Database = Mongo.Client.Database(config.DB_NAME)

	log.Printf("Conectado exitosamente a MongoDB: %s - Base de datos: %s", config.MONGO_URI, config.DB_NAME)
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
