package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
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
