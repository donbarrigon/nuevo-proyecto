package model

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Permission struct {
	ID   bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string        `bson:"name"          json:"name"`
}

func (p *Permission) CollectionName() string { return "permissions" }

func (p *Permission) GetID() bson.ObjectID { return p.ID }

func (p *Permission) SetID(id bson.ObjectID) { p.ID = id }

func (p *Permission) BeforeCreate() app.Error { return nil }

func (p *Permission) BeforeUpdate() app.Error { return nil }
