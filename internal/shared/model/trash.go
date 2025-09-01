package model

import (
	"context"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	. "github.com/donbarrigon/nuevo-proyecto/internal/app/qb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Trash struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Collection string        `bson:"collection"    json:"collection"`
	Document   any           `bson:"document"      json:"document"`
	DeletedAt  time.Time     `bson:"deleted_at"    json:"deleted_at"`
	app.Odm    `bson:"-" json:"-"`
}

func NewTrash() *Trash {
	trash := &Trash{}
	trash.Odm.Model = trash
	return trash
}

func (t *Trash) CollectionName() string { return "trash" }
func (t *Trash) GetID() bson.ObjectID   { return t.ID }
func (t *Trash) SetID(id bson.ObjectID) { t.ID = id }

func (t *Trash) BeforeCreate() app.Error {
	t.DeletedAt = time.Now()
	return nil
}

func (t *Trash) BeforeUpdate() app.Error {
	return app.Errors.Unknownf("you tried to modify an trash record")
}

func (t *Trash) MoveToTrash(m app.Model) app.Error {
	t.Collection = m.CollectionName()
	t.Document = m
	if err := t.Create(); err != nil {
		return err
	}
	return m.Delete()
}

func (t *Trash) RestoreByHex(m app.Model, id string) app.Error {
	oid, er := bson.ObjectIDFromHex(id)
	if er != nil {
		return app.Errors.HexID(er)
	}
	return t.Restore(m, oid)
}

func (t *Trash) Restore(m app.Model, id bson.ObjectID) app.Error {
	ctx := context.TODO()
	cursor, er := app.DB.Collection(t.CollectionName()).Aggregate(ctx, Pipeline(
		Match(Where("document._id", Eq(id))),
		bson.D{{"$sort", bson.D{{"deleted_at", -1}}}},
		bson.D{{"$limit", 1}},
		bson.D{{"$replaceRoot", bson.D{{"newRoot", "$document"}}}},
	))
	if er != nil {
		return app.Errors.Mongo(er)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if er := cursor.Decode(m); er != nil {
			return app.Errors.Mongo(er)
		}
	} else {
		return app.Errors.NoDocumentsf("No documents matched for restore :collection [:model::id]", app.E("id", id), app.E("model", m.CollectionName()), app.E("collection", "activity"))
	}
	if err := m.Create(); err != nil {
		return err
	}
	return t.DeleteOne(Filter(Where("document._id", Eq(id))))
}

func (t *Trash) PurgeByHex(id string) app.Error {
	oid, er := bson.ObjectIDFromHex(id)
	if er != nil {
		return app.Errors.HexID(er)
	}
	return t.Purge(oid)
}

func (t *Trash) Purge(id bson.ObjectID) app.Error {
	return t.DeleteOne(Filter(Where("document._id", Eq(id))))
}
