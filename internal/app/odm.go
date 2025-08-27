package app

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Model interface {
	CollectionName() string
	GetID() bson.ObjectID
	SetID(id bson.ObjectID)
	BeforeCreate() Error
	BeforeUpdate() Error
}

type Collection []Model

type Odm struct {
	Model Model `bson:"-" json:"-"`
}

var DBClient *mongo.Client
var DB *mongo.Database

// trae el documento segun el id en string
func (o *Odm) FindByHexID(id string) Error {

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return Errors.HexID(err)
	}
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	if err := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); err != nil {
		return Errors.Mongo(err)
	}
	return nil
}

// trae el documento segun el id
func (o *Odm) FindByID(id bson.ObjectID) Error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	if err := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); err != nil {
		return Errors.Mongo(err)
	}
	return nil
}

// trae el primer documento encontrado
func (o *Odm) First(field string, value any) Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	if err := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); err != nil {
		return Errors.Mongo(err)
	}
	return nil
}

// trae 1 documento segun el filtro
func (o *Odm) FindOne(filter bson.D, opts ...options.Lister[options.FindOneOptions]) Error {
	if err := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter, opts...).Decode(o.Model); err != nil {
		return Errors.Mongo(err)
	}
	return nil
}

// ejecuta busquedas por el filtro
func (o *Odm) Find(result any, filter bson.D, opts ...options.Lister[options.FindOptions]) Error {
	ctx := context.TODO()
	cursor, err := DB.Collection(o.Model.CollectionName()).Find(ctx, filter)
	if err != nil {
		return Errors.Mongo(err)
	}
	if err = cursor.All(ctx, result); err != nil {
		return Errors.Mongo(err)
	}
	return nil
}

// trae un slice con todos los documentos encontrados
func (o *Odm) FindBy(result any, field string, value any) Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	ctx := context.TODO()
	cursor, err := DB.Collection(o.Model.CollectionName()).Find(ctx, filter)
	if err != nil {
		return Errors.Mongo(err)
	}
	if err = cursor.All(ctx, result); err != nil {
		return Errors.Mongo(err)
	}
	return nil
}

// ejecuta busquedas por agregacion
func (o *Odm) Aggregate(result any, pipeline mongo.Pipeline) Error {
	ctx := context.TODO()
	cursor, err := DB.Collection(o.Model.CollectionName()).Aggregate(ctx, pipeline)
	if err != nil {
		return Errors.Mongo(err)
	}
	if err = cursor.All(ctx, result); err != nil {
		return Errors.Mongo(err)
	}
	return nil
}

// lo mismo que agregate pero solo retorna 1
func (o *Odm) AggregateOne(pipeline mongo.Pipeline) Error {
	ctx := context.TODO()
	cursor, err := DB.Collection(o.Model.CollectionName()).Aggregate(ctx, pipeline)
	if err != nil {
		return Errors.Mongo(err)
	}
	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		if err := cursor.Decode(o.Model); err != nil {
			return Errors.Mongo(err)
		}
	}
	return nil
}

// crea el documento
func (o *Odm) Create() Error {
	if err := o.Model.BeforeCreate(); err != nil {
		return err
	}
	result, err := DB.Collection(o.Model.CollectionName()).InsertOne(context.TODO(), o.Model)
	if err != nil {
		return Errors.Mongo(err)
	}
	o.Model.SetID(result.InsertedID.(bson.ObjectID))

	return nil
}

// crea el documento con los datos del validador
func (o *Odm) CreateBy(validator any) Error {
	if err := Fill(o.Model, validator); err != nil {
		return err
	}
	return o.Create()
}

// crea varios documentos los datos deben ser un slice del modelo
func (o *Odm) CreateMany(data any) Error {

	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return Errors.InternalServerErrorf("Create many required a slice")
	}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		if err := elem.(Model).BeforeCreate(); err != nil {
			return err
		}
	}
	collection := DB.Collection(o.Model.CollectionName())
	result, err := collection.InsertMany(context.TODO(), data)
	if err != nil {
		return Errors.Mongo(err)
	}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		elem.(Model).SetID(result.InsertedIDs[i].(bson.ObjectID))
	}
	return nil
}

// actualiza el documento
func (o *Odm) Update() Error {
	if err := o.Model.BeforeUpdate(); err != nil {
		return err
	}
	filter := bson.D{bson.E{Key: "_id", Value: o.Model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: o.Model}}

	result, err := DB.Collection(o.Model.CollectionName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Errors.Mongo(err)
	}
	if result.MatchedCount == 0 {
		return Errors.NoDocumentsf("mongo.UpdateResult.MatchedCount == 0")
	}
	if result.ModifiedCount == 0 {
		return Errors.Updatef("mongo.UpdateResult.ModifiedCount == 0")
	}
	return nil
}

// actualiza el documento con los datos del validador
func (o *Odm) UpdateBy(validator any) (map[string]any, Error) {
	dirty, err := FillDirty(o.Model, validator)
	if err != nil {
		return nil, err
	}
	return dirty, o.Update()
}

func (o *Odm) UpdateOne(filter bson.D, update bson.D) Error {
	if err := o.Model.BeforeUpdate(); err != nil {
		return err
	}
	result, err := DB.Collection(o.Model.CollectionName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Errors.Mongo(err)
	}
	if result.MatchedCount == 0 {
		return Errors.NoDocumentsf("mongo.UpdateResult.MatchedCount == 0")
	}
	if result.ModifiedCount == 0 {
		return Errors.Updatef("mongo.UpdateResult.ModifiedCount == 0")
	}
	return nil
}

// hace un soft delete al documento
func (o *Odm) SoftDelete() Error {
	collection := DB.Collection(o.Model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: o.Model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Errors.Mongo(err)
	}
	if result.MatchedCount == 0 {
		return Errors.NoDocumentsf("mongo.UpdateResult.MatchedCount == 0")
	}
	if result.ModifiedCount == 0 {
		return Errors.Deletef("mongo.UpdateResult.ModifiedCount == 0")
	}
	return nil
}

// restaura el documento eliminado por SoftDelete
func (o *Odm) Restore() Error {
	filter := bson.D{bson.E{Key: "_id", Value: o.Model.GetID()}}
	update := bson.D{bson.E{Key: "$unset", Value: bson.D{{Key: "deleted_at", Value: nil}}}}

	result, err := DB.Collection(o.Model.CollectionName()).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Errors.Mongo(err)
	}
	if result.MatchedCount == 0 {
		return Errors.NoDocumentsf("mongo.UpdateResult.MatchedCount == 0")
	}
	if result.ModifiedCount == 0 {
		return Errors.Restoref("mongo.UpdateResult.ModifiedCount == 0")
	}
	return nil
}

// elimina permanentemente el documento
func (o *Odm) Delete() Error {
	collection := DB.Collection(o.Model.CollectionName())
	filter := bson.D{bson.E{Key: "_id", Value: o.Model.GetID()}}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return Errors.Mongo(err)
	}
	if result.DeletedCount == 0 {
		return Errors.ForceDeletef("mongo.DeleteResult.DeletedCount == 0")
	}
	return nil
}
