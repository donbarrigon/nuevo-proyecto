package db

import (
	"context"
	goerrors "errors"
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

func (db *ConexionMongoDB) FindByHexID(model MongoModel, id string) errors.Error {

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.HexID(err)
	}
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	if err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func (db *ConexionMongoDB) FindByID(model MongoModel, id bson.ObjectID) errors.Error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	if err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func (db *ConexionMongoDB) FindManyByField(model MongoModel, result any, field string, value any) errors.Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	cursor, err := db.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func (db *ConexionMongoDB) FindOneByField(model MongoModel, field string, value any) errors.Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	if err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func (db *ConexionMongoDB) FindAll(model MongoModel, result any) errors.Error {

	cursor, err := db.Database.Collection(model.CollectionName()).Find(context.TODO(), bson.D{})
	if err != nil {
		return errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return errors.Mongo(err)
	}
	return nil

}

func (db *ConexionMongoDB) FindOne(model MongoModel, filter bson.D) errors.Error {
	err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model)
	if err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func (db *ConexionMongoDB) FindMany(model MongoModel, result any, filter bson.D) errors.Error {
	cursor, err := db.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func (db *ConexionMongoDB) Create(model MongoModel) errors.Error {
	model.Default()
	collection := db.Database.Collection(model.CollectionName())
	result, err := collection.InsertOne(context.TODO(), model)
	if err != nil {
		return errors.Mongo(err)
	}
	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return errors.Unknown(goerrors.New("No se logro hacer la conversion a bson.ObjectID"))
	}
	model.SetID(id)
	return nil
}

func (db *ConexionMongoDB) Update(model MongoModel) (*mongo.UpdateResult, errors.Error) {
	model.Default()
	collection := db.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: model}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return result, errors.Mongo(err)
	}
	return result, nil
}
func (db *ConexionMongoDB) Delete(model MongoModel) (*mongo.UpdateResult, errors.Error) {
	collection := db.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{{Key: "deletedAt", Value: time.Now()}}}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return result, errors.Mongo(err)
	}
	return result, nil
}

func (db *ConexionMongoDB) Restore(model MongoModel) (*mongo.UpdateResult, errors.Error) {
	collection := db.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$unset", Value: bson.D{{Key: "deletedAt", Value: nil}}}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return result, errors.Mongo(err)
	}
	return result, nil
}

func (db *ConexionMongoDB) ForceDelete(model MongoModel) (*mongo.DeleteResult, errors.Error) {
	collection := db.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return result, errors.Mongo(err)
	}
	return result, nil
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
