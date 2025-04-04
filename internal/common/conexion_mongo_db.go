package com

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

func (db *ConexionMongoDB) FindByHexID(model Model, id string) Error {

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return &Err{
			Status:  http.StatusBadRequest,
			Message: TT(LANG, "El id [%v] no es un hexadecimal válido", id),
			Err:     err.Error(),
		}
	}
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	if err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		if err == mongo.ErrNoDocuments {
			return &Err{
				Status:  http.StatusNotFound,
				Message: TT(LANG, "No se encontró el registro [%v]", id),
				Err:     err.Error(),
			}
		}
		return &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al buscar el registro [%v]", id),
			Err:     err.Error(),
		}
	}
	return nil
}

func (db *ConexionMongoDB) FindByID(model Model, id bson.ObjectID) Error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	if err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		if err == mongo.ErrNoDocuments {
			return &Err{
				Status:  http.StatusNotFound,
				Message: TT(LANG, "No se encontró el registro [%v]", id.String()),
				Err:     err.Error(),
			}
		}
		return &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al buscar el registro [%v]", id.String()),
			Err:     err.Error(),
		}
	}
	return nil
}

func (db *ConexionMongoDB) FindManyByField(model Model, field string, value any) (*[]Model, Error) {
	filter := bson.D{bson.E{Key: field, Value: value}}
	cursor, err := db.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &Err{
				Status:  http.StatusNotFound,
				Message: TT(LANG, "No se encontraron registros con el campo [%v]: %v", field, value),
				Err:     err.Error(),
			}
		}
		return nil, &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al buscar el registro por el campo [%v]: %v", field, value),
			Err:     err.Error(),
		}
	}
	results := toSlice(model)
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al recuperar los resultados de la búsqueda por el campo [%v]: %v", field, value),
			Err:     err.Error(),
		}
	}
	return &results, nil
}

func (db *ConexionMongoDB) FindOneByField(model Model, field string, value any) Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	if err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		if err == mongo.ErrNoDocuments {
			return &Err{
				Status:  http.StatusNotFound,
				Message: TT(LANG, "No se encontró el registro con el campo [%v]: %v", field, value),
				Err:     err.Error(),
			}
		}
		return &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al buscar el registro con el campo [%v]: %v", field, value),
			Err:     err.Error(),
		}
	}
	return nil
}

func (db *ConexionMongoDB) FindOne(model Model, filter bson.D) Error {
	err := db.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &Err{
				Status:  http.StatusNotFound,
				Message: TT(LANG, "No se encontró el registro"),
				Err:     err.Error(),
			}
		}
		return &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al buscar el registro"),
			Err:     err.Error(),
		}
	}
	return nil
}

func (db *ConexionMongoDB) FindMany(model Model, filter bson.D) (*[]Model, Error) {
	cursor, err := db.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &Err{
				Status:  http.StatusNotFound,
				Message: TT(LANG, "No se encontraron registros"),
				Err:     err.Error(),
			}
		}
		return nil, &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al buscar los registros"),
			Err:     err.Error(),
		}
	}
	results := toSlice(model)
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al recuperar los resultados de la búsqueda"),
			Err:     err.Error(),
		}
	}
	return &results, nil
}

func toSlice[T any](v T) []T {
	return make([]T, 0)
}

func (db *ConexionMongoDB) Create(model Model) (*mongo.InsertOneResult, Error) {
	collection := db.Database.Collection(model.CollectionName())
	result, err := collection.InsertOne(context.TODO(), model)
	if err != nil {
		return nil, &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al insertar el registro"),
			Err:     err.Error(),
		}
	}
	return result, nil
}

func (db *ConexionMongoDB) Update(model Model) (*mongo.UpdateResult, Error) {
	collection := db.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: model}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al actualizar el registro"),
			Err:     err.Error(),
		}
	}
	return result, nil
}
func (db *ConexionMongoDB) Delete(model Model) (*mongo.UpdateResult, Error) {
	collection := db.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{{Key: "deletedAt", Value: time.Now()}}}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al eliminar el registro"),
			Err:     err.Error(),
		}
	}
	return result, nil
}

func (db *ConexionMongoDB) Restore(model Model) (*mongo.UpdateResult, Error) {
	collection := db.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$unset", Value: bson.D{{Key: "deletedAt", Value: nil}}}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al restaurar el registro"),
			Err:     err.Error(),
		}
	}
	return result, nil
}

func (db *ConexionMongoDB) ForceDelete(model Model) (*mongo.DeleteResult, Error) {
	collection := db.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, &Err{
			Status:  http.StatusInternalServerError,
			Message: TT(LANG, "Error al eliminar permanentemente el registro"),
			Err:     err.Error(),
		}
	}
	return result, nil
}

func InitMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(MONGO_URI).SetServerAPIOptions(serverAPI)
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
