package service

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/shared/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ActivityRecord(userID bson.ObjectID, collection app.Model, action string, changes ...any) {
	if len(changes) == 0 {
		changes = append(changes, map[string]any{})
	}
	activity := &model.Activity{
		UserID:     userID,
		DocumentID: collection.GetID(),
		Collection: collection.CollectionName(),
		Action:     action,
		Changes:    changes,
	}
	activity.Odm.Model = activity
	if err := activity.Create(); err != nil {
		app.PrintError("Failed to create activity record", app.Entry{Key: "error", Value: err})
	}
}

func ActivityManyRecords(userID bson.ObjectID, collection app.Model, action string, changes ...any) {
	activity := &model.Activity{}
	activity.Odm.Model = activity
	data := make([]*model.Activity, len(changes))
	for _, change := range changes {
		if change == nil {
			change = bson.M{}
		}
		data = append(data, &model.Activity{
			UserID:     userID,
			DocumentID: collection.GetID(),
			Collection: collection.CollectionName(),
			Action:     action,
			Changes:    change,
		})
	}
	if err := activity.CreateMany(data); err != nil {
		app.PrintError("Failed to create activity record", app.Entry{Key: "error", Value: err})
	}
}
