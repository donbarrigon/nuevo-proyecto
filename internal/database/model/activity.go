package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Activity struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       string        `bson:"user_id"       json:"user_id"`
	CollectionID bson.ObjectID `bson:"collection_id" json:"collection_id"`
	Collection   string        `bson:"collection"    json:"collection"`
	Action       string        `bson:"action"        json:"action"`
	Changes      any           `bson:"changes"       json:"changes"`
	CreatedAt    time.Time     `bson:"created_at"    json:"created_at"`
	app.Orm
}

func NewActivity() *Activity {
	activity := &Activity{}
	activity.Orm.Model = activity
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

func ActivityRecord(model app.Model, action string, changes any) {
	if changes == nil {
		changes = bson.M{}
	}
	activity := &Activity{
		CollectionID: model.GetID(),
		Collection:   model.CollectionName(),
		Action:       action,
		Changes:      changes,
	}
	activity.Orm.Model = activity
	if err := activity.Create(); err != nil {
		app.Log.Error("Failed to create activity record", app.F{Key: "error", Value: err})
	}
}

func ActivityManyRecords(model app.Model, action string, changes []any) {
	activity := &Activity{}
	activity.Orm.Model = activity
	data := make([]*Activity, len(changes))
	for _, change := range changes {
		if change == nil {
			change = bson.M{}
		}
		data = append(data, &Activity{
			CollectionID: model.GetID(),
			Collection:   model.CollectionName(),
			Action:       action,
			Changes:      change,
		})
	}
	if err := activity.CreateMany(data); err != nil {
		app.Log.Error("Failed to create activity record", app.F{Key: "error", Value: err})
	}
}
