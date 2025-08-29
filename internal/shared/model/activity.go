package model

import (
	"context"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	. "github.com/donbarrigon/nuevo-proyecto/internal/app/qb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Activity struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     bson.ObjectID `bson:"user_id"       json:"user_id"`
	DocumentID bson.ObjectID `bson:"document_id" json:"document_id"`
	Collection string        `bson:"collection"    json:"collection"`
	Action     string        `bson:"action"        json:"action"`
	Changes    any           `bson:"changes"       json:"changes"`
	CreatedAt  time.Time     `bson:"created_at"    json:"created_at"`
	app.Odm    `bson:"-" json:"-"`
}

func NewActivity() *Activity {
	activity := &Activity{}
	activity.Odm.Model = activity
	return activity
}

func (a *Activity) CollectionName() string { return "activities" }
func (a *Activity) GetID() bson.ObjectID   { return a.ID }
func (a *Activity) SetID(id bson.ObjectID) { a.ID = id }

func (a *Activity) BeforeCreate() app.Error {
	a.CreatedAt = time.Now()
	return nil
}

func (a *Activity) BeforeUpdate() app.Error {
	return app.Errors.Unknownf("you tried to modify an activity record")
}

func (a *Activity) RecoverDocument(m app.Model, id string) app.Error {
	ctx := context.TODO()
	oid, er := bson.ObjectIDFromHex(id)
	if er != nil {
		return app.Errors.HexID(er)
	}
	cursor, er := app.DB.Collection(m.CollectionName()).Aggregate(ctx, Pipeline(
		Match(
			Where("document_id", Eq(oid)),
			Where("collection", Eq(m.CollectionName())),
			Where("action", Eq("delete")),
		),
		bson.D{{"$sort", bson.D{{"created_at", -1}}}},
		bson.D{{"$limit", 1}},
		bson.D{{"$replaceRoot", bson.D{{"newRoot", "$changes"}}}},
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
	return nil
}
