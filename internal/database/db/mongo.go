package db

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
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
	Default()
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

func Find(model MongoModel, result any, filter bson.D) errors.Error {
	cursor, err := Mongo.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return errors.Mongo(err)
	}
	return nil
}

func FindByPipeline(model MongoModel, result any, pipeline any) errors.Error {
	cursor, err := Mongo.Database.Collection(model.CollectionName()).Aggregate(context.TODO(), pipeline)
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

	if err := setID(model, result.InsertedID); err != nil {
		return errors.Unknown(err)
	}

	return nil
}

func Update(model MongoModel) errors.Error {
	model.Default()
	collection := Mongo.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: getID(model)}}
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
	filter := bson.D{bson.E{Key: "_id", Value: getID(model)}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}}}

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
	filter := bson.D{bson.E{Key: "_id", Value: getID(model)}}
	update := bson.D{bson.E{Key: "$unset", Value: bson.D{{Key: "deleted_at", Value: nil}}}}

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
	filter := bson.D{bson.E{Key: "_id", Value: getID(model)}}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return errors.Mongo(err)
	}

	if result.DeletedCount == 0 {
		return errors.ForceDelete(errors.New("mongo.DeleteResult.DeletedCount == 0"))
	}
	return nil
}

func setID(model MongoModel, id any) error {
	val := reflect.ValueOf(model)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("el modelo debe ser un puntero no nulo")
	}
	val = val.Elem()

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("bson")
		if strings.Split(tag, ",")[0] == "_id" {
			fieldVal := val.Field(i)
			if !fieldVal.CanSet() {
				return fmt.Errorf("no se puede asignar el campo _id")
			}

			idVal := reflect.ValueOf(id)

			// Si el tipo es exactamente igual
			if fieldVal.Type() == idVal.Type() {
				fieldVal.Set(idVal)
				return nil
			}

			// Si el campo es string y el id es bson.ObjectID
			if fieldVal.Kind() == reflect.String && idVal.Type().String() == "primitive.ObjectID" {
				method := idVal.MethodByName("Hex")
				if method.IsValid() && method.Type().NumIn() == 0 {
					hexVal := method.Call(nil)[0]
					fieldVal.SetString(hexVal.String())
					return nil
				}
			}

			return fmt.Errorf("no se pudo asignar el ID: tipos incompatibles")
		}
	}
	return fmt.Errorf("no se encontró el campo con tag bson:\"_id\"")
}

func getID(model MongoModel) any {
	val := reflect.ValueOf(model)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("bson")
		tagParts := strings.Split(tag, ",")
		if len(tagParts) > 0 && tagParts[0] == "_id" {
			return val.Field(i).Interface()
		}
	}

	return nil
}

func InitMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(config.Env.DB_URI).SetServerAPIOptions(serverAPI)
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
	Mongo.Database = Mongo.Client.Database(config.Env.DB_DATABASE)

	log.Printf("Conectado exitosamente a MongoDB: %s - Base de datos: %s", config.Env.DB_URI, config.Env.DB_DATABASE)
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
