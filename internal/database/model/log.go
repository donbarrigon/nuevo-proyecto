package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Log struct {
	ID       bson.ObjectID     `bson:"id,omitempty" `
	Time     string            `bson:"time,omitempty"`
	Level    string            `bson:"level,omitempty"`
	Message  string            `bson:"message"`
	Function string            `bson:"function,omitempty" `
	Line     string            `bson:"line,omitempty"`
	File     string            `bson:"file,omitempty"`
	Context  map[string]string `bson:"context,omitempty"`
}

func NewLog() *Log {
	return &Log{
		ID: bson.NewObjectID(),
	}
}

func (l *Log) TableName() string {
	return "permissions"
}

func (l *Log) Default() {
	//
}
