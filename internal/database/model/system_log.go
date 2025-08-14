package model

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SystemLog struct {
	ID       bson.ObjectID     `bson:"id,omitempty"       json:"id"`
	Time     string            `bson:"time,omitempty"     json:"time"`
	Level    string            `bson:"level,omitempty"    json:"level"`
	Message  string            `bson:"message"            json:"message"`
	Function string            `bson:"function,omitempty" json:"function,omitempty"`
	Line     string            `bson:"line,omitempty"     json:"line,omitempty"`
	File     string            `bson:"file,omitempty"     json:"file,omitempty"`
	Context  map[string]string `bson:"context,omitempty"  json:"context,omitempty"`
}

func (l *SystemLog) CollectionName() string { return "system_logs" }

func (l *SystemLog) GetID() bson.ObjectID { return l.ID }

func (l *SystemLog) SetID(id bson.ObjectID) { l.ID = id }

func (l *SystemLog) BeforeCreate() app.Error { return nil }

func (l *SystemLog) BeforeUpdate() app.Error { return nil }
