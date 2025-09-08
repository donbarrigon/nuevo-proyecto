package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SystemLog struct {
	ID       bson.ObjectID     `bson:"id,omitempty"       json:"id"`
	Time     time.Time         `bson:"time,omitempty"     json:"time"`
	Level    string            `bson:"level,omitempty"    json:"level"`
	Message  string            `bson:"message"            json:"message"`
	Function string            `bson:"function,omitempty" json:"function,omitempty"`
	Line     string            `bson:"line,omitempty"     json:"line,omitempty"`
	File     string            `bson:"file,omitempty"     json:"file,omitempty"`
	Context  map[string]string `bson:"context,omitempty"  json:"context,omitempty"`
	app.Odm  `bson:"-" json:"-"`
}

func NewSystemLog() *SystemLog {
	systemLog := &SystemLog{}
	systemLog.Odm.Model = systemLog
	return systemLog
}

func (l *SystemLog) CollectionName() string { return "system_logs" }
func (l *SystemLog) GetID() bson.ObjectID   { return l.ID }
func (l *SystemLog) SetID(id bson.ObjectID) { l.ID = id }

func (l *SystemLog) BeforeCreate() app.Error { return nil }

func (l *SystemLog) BeforeUpdate() app.Error { return nil }
