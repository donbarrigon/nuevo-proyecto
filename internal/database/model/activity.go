package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Activity struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string        `bson:"user_id"       json:"user_id"`
	ModelID   bson.ObjectID `bson:"model_id"      json:"model_id"`
	Model     string        `bson:"model"         json:"model"`
	Action    string        `bson:"action"        json:"action"`
	Changes   any           `bson:"changes"       json:"changes"`
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`
}

func (a *Activity) CollectionName() string { return "activities" }

func (a *Activity) GetID() bson.ObjectID { return a.ID }

func (a *Activity) SetID(id bson.ObjectID) { a.ID = id }

func (a *Activity) BeforeCreate() app.Error {
	a.CreatedAt = time.Now()
	return nil
}

func (a *Activity) BeforeUpdate() app.Error {
	return app.Errors.Unknownf("you tried to modify an activity record")
}

func ActivityRecord(model db.Model, action string, changes any) {
	if changes == nil {
		changes = bson.M{}
	}
	activity := &Activity{
		ModelID: model.GetID(),
		Model:   model.CollectionName(),
		Action:  action,
		Changes: changes,
	}
	if err := db.Create(activity); err != nil {
		app.Log.Error("Failed to create activity record", app.F{Key: "error", Value: err})
	}
}

func ActivityManyRecords(model db.Model, action string, changes []any) {
	data := make([]*Activity, len(changes))
	for _, change := range changes {
		if change == nil {
			change = bson.M{}
		}
		data = append(data, &Activity{
			ModelID: model.GetID(),
			Model:   model.CollectionName(),
			Action:  action,
			Changes: change,
		})
	}
	if err := db.CreateMany(&Activity{}, data); err != nil {
		app.Log.Error("Failed to create activity record", app.F{Key: "error", Value: err})
	}
}
