package model

import (
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Permission struct {
	ID   bson.ObjectID `bson:"_id,omitempty" json:"-"`
	Name string        `bson:"name" json:"name"`
}

func NewPermission() *Permission {
	return &Permission{
		ID: bson.NewObjectID(),
	}
}
func (p *Permission) CollectionName() string {
	return "permissions"
}

func (p *Permission) Default() {
	//...
}

func (p *Permission) Validate(l string) errors.Error {
	err := &errors.Err{}

	if len(p.Name) > 255 {
		err.Append("name", lang.TT(l, "Maximo %v caracteres", 255))
	}

	if p.Name == "" {
		err.Append("name", lang.TT(l, "Este campo es requerido"))
	}
	return err.Errors()
}
