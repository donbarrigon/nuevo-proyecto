package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type History struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     bson.ObjectID `bson:"user_id"       json:"user_id"`
	DocumentID bson.ObjectID `bson:"document_id"   json:"document_id"`
	Collection string        `bson:"collection"    json:"collection"`
	Action     string        `bson:"action"        json:"action"`
	Old        any           `bson:"old"           json:"old"`
	OcurredAt  time.Time     `bson:"occurred_at"   json:"occurred_at"`
	app.Odm    `bson:"-" json:"-"`
}

const (
	ACTION_CREATE        = "create"
	ACTION_UPDATE        = "update"
	ACTION_DELETE        = "delete"
	ACTION_MOVE_TO_TRASH = "move-to-trash"
	ACTION_RESTORE       = "restore"
	ACTION_PURGUE        = "purge"
)

func NewHistory() *History {
	history := &History{}
	history.Odm.Model = history
	return history
}

func (a *History) CollectionName() string { return "histories" }
func (a *History) GetID() bson.ObjectID   { return a.ID }
func (a *History) SetID(id bson.ObjectID) { a.ID = id }

func (a *History) BeforeCreate() app.Error {
	a.OcurredAt = time.Now()
	return nil
}

func (a *History) BeforeUpdate() app.Error {
	return app.Errors.Unknownf("you tried to modify an history record")
}

func HistoryRecord(userID bson.ObjectID, collection app.Model, action string, old any) {
	if old == nil {
		old = map[string]any{}
	}
	history := &History{
		UserID:     userID,
		DocumentID: collection.GetID(),
		Collection: collection.CollectionName(),
		Action:     action,
		Old:        old,
	}
	history.Odm.Model = history
	if err := history.Create(); err != nil {
		app.PrintError("Failed to create hiistory record", app.Entry{Key: "error", Value: err})
	}
}

func HistoryManyRecords(userID bson.ObjectID, collection app.Model, action string, old ...any) {
	history := &History{}
	history.Odm.Model = history
	data := make([]History, len(old))
	for _, change := range old {
		if change == nil {
			change = bson.M{}
		}
		data = append(data, History{
			UserID:     userID,
			DocumentID: collection.GetID(),
			Collection: collection.CollectionName(),
			Action:     action,
			Old:        change,
		})
	}
	if err := history.CreateMany(data); err != nil {
		app.PrintError("Failed to create activity record", app.Entry{Key: "error", Value: err})
	}
}
