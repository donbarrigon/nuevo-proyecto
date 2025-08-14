package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Model interface {
	CollectionName() string
	GetID() bson.ObjectID
	SetID(id bson.ObjectID)
	BeforeCreate() app.Error
	BeforeUpdate() app.Error
}

type ConexionMongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var Mongo *ConexionMongoDB

func FindByHexID(model Model, id string) app.Error {

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return app.Errors.HexID(err)
	}
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	if err := Mongo.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		return app.Errors.Mongo(err)
	}
	return nil
}

func FindByID(model Model, id bson.ObjectID) app.Error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	if err := Mongo.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		return app.Errors.Mongo(err)
	}
	return nil
}

func FindBy(model Model, result any, field string, value any) app.Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	cursor, err := Mongo.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return app.Errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return app.Errors.Mongo(err)
	}
	return nil
}

func FindOneBy(model Model, field string, value any) app.Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	if err := Mongo.Database.Collection(model.CollectionName()).FindOne(context.TODO(), filter).Decode(model); err != nil {
		return app.Errors.Mongo(err)
	}
	return nil
}

func Find(model Model, result any, filter bson.D) app.Error {
	cursor, err := Mongo.Database.Collection(model.CollectionName()).Find(context.TODO(), filter)
	if err != nil {
		return app.Errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return app.Errors.Mongo(err)
	}
	return nil
}

func Aggregate(model Model, result any, pipeline mongo.Pipeline) app.Error {
	cursor, err := Mongo.Database.Collection(model.CollectionName()).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return app.Errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return app.Errors.Mongo(err)
	}
	return nil
}

func Create(model Model) app.Error {
	if err := model.BeforeCreate(); err != nil {
		return err
	}
	collection := Mongo.Database.Collection(model.CollectionName())
	result, err := collection.InsertOne(context.TODO(), model)
	if err != nil {
		return app.Errors.Mongo(err)
	}
	model.SetID(result.InsertedID.(bson.ObjectID))

	return nil
}

func CreateBy(model Model, validator any) app.Error {
	if err := Fill(model, validator); err != nil {
		return err
	}
	return Create(model)
}

func CreateMany(model Model, data any) app.Error {
	models, ok := data.([]Model)
	if !ok {
		return app.Errors.InternalServerErrorF("type assertion failed in CreateMany")
	}
	for _, m := range models {

		if err := m.BeforeCreate(); err != nil {
			return err
		}
	}
	collection := Mongo.Database.Collection(model.CollectionName())
	result, err := collection.InsertMany(context.TODO(), data)
	if err != nil {
		return app.Errors.Mongo(err)
	}
	for i, m := range models {
		m.SetID(result.InsertedIDs[i].(bson.ObjectID))
	}
	return nil
}

func Update(model Model) app.Error {
	if err := model.BeforeUpdate(); err != nil {
		return err
	}
	collection := Mongo.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: model}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return app.Errors.Mongo(err)
	}
	if result.MatchedCount == 0 {
		return app.Errors.NoDocumentsf("mongo.UpdateResult.MatchedCount == 0")
	}
	if result.ModifiedCount == 0 {
		return app.Errors.Updatef("mongo.UpdateResult.ModifiedCount == 0")
	}
	return nil
}

func UpdateBy(model Model, validator any) (map[string]any, app.Error) {
	dirty, err := FillDirty(model, validator)
	if err != nil {
		return nil, err
	}
	return dirty, Update(model)
}

func Delete(model Model) app.Error {
	collection := Mongo.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return app.Errors.Mongo(err)
	}
	if result.MatchedCount == 0 {
		return app.Errors.NoDocumentsf("mongo.UpdateResult.MatchedCount == 0")
	}
	if result.ModifiedCount == 0 {
		return app.Errors.Deletef("mongo.UpdateResult.ModifiedCount == 0")
	}
	return nil
}

func Restore(model Model) app.Error {
	collection := Mongo.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}
	update := bson.D{bson.E{Key: "$unset", Value: bson.D{{Key: "deleted_at", Value: nil}}}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return app.Errors.Mongo(err)
	}
	if result.MatchedCount == 0 {
		return app.Errors.NoDocumentsf("mongo.UpdateResult.MatchedCount == 0")
	}
	if result.ModifiedCount == 0 {
		return app.Errors.Restoref("mongo.UpdateResult.ModifiedCount == 0")
	}
	return nil
}

func ForceDelete(model Model) app.Error {
	collection := Mongo.Database.Collection(model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: model.GetID()}}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return app.Errors.Mongo(err)
	}
	if result.DeletedCount == 0 {
		return app.Errors.ForceDeletef("mongo.DeleteResult.DeletedCount == 0")
	}
	return nil
}

func InitMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(app.Env.DB_CONNECTION_STRING).SetServerAPIOptions(serverAPI)
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
	Mongo.Database = Mongo.Client.Database(app.Env.DB_DATABASE)

	app.Log.Print("Conectado exitosamente a MongoDB: :con - Base de datos: :db",
		app.F{Key: "con", Value: app.Env.DB_CONNECTION_STRING},
		app.F{Key: "con", Value: app.Env.DB_DATABASE})
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

// deprecado por que consume muchos recursos
// func setID(model Model, id any) error {
// 	val := reflect.ValueOf(model)
// 	if val.Kind() != reflect.Ptr || val.IsNil() {
// 		return fmt.Errorf("el modelo debe ser un puntero no nulo")
// 	}
// 	val = val.Elem()

// 	typ := val.Type()
// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		tag := field.Tag.Get("bson")
// 		if strings.Split(tag, ",")[0] == "_id" {
// 			fieldVal := val.Field(i)
// 			if !fieldVal.CanSet() {
// 				return fmt.Errorf("no se puede asignar el campo _id")
// 			}

// 			idVal := reflect.ValueOf(id)

// 			// Si el tipo es exactamente igual
// 			if fieldVal.Type() == idVal.Type() {
// 				fieldVal.Set(idVal)
// 				return nil
// 			}

// 			// Si el campo es string y el id es bson.ObjectID
// 			if fieldVal.Kind() == reflect.String && idVal.Type().String() == "primitive.ObjectID" {
// 				method := idVal.MethodByName("Hex")
// 				if method.IsValid() && method.Type().NumIn() == 0 {
// 					hexVal := method.Call(nil)[0]
// 					fieldVal.SetString(hexVal.String())
// 					return nil
// 				}
// 			}

// 			return fmt.Errorf("no se pudo asignar el ID: tipos incompatibles")
// 		}
// 	}
// 	return fmt.Errorf("no se encontró el campo con tag bson:\"_id\"")
// }

// func getID(model Model) any {
// 	val := reflect.ValueOf(model)

// 	if val.Kind() == reflect.Ptr {
// 		if val.IsNil() {
// 			return nil
// 		}
// 		val = val.Elem()
// 	}

// 	if val.Kind() != reflect.Struct {
// 		return nil
// 	}

// 	typ := val.Type()

// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		tag := field.Tag.Get("bson")
// 		tagParts := strings.Split(tag, ",")
// 		if len(tagParts) > 0 && tagParts[0] == "_id" {
// 			return val.Field(i).Interface()
// 		}
// 	}

// 	return nil
// }
