package model

import (
	"strings"

	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"-"`
	Name        string        `bson:"name" json:"name"`
	Permissions []*Permission `bson:"permissions" json:"-"`
}

func NewRole() *Role {
	return &Role{
		ID: bson.NewObjectID(),
	}
}
func (r *Role) CollectionName() string {
	return "roles"
}
func (r *Role) GetID() bson.ObjectID {
	return r.ID
}

func (r *Role) SetID(id bson.ObjectID) {
	r.ID = id
}

func (r *Role) Default() {
	//...
}

func (r *Role) Validate(l string) errors.Error {
	err := &errors.Err{}

	if strings.TrimSpace(r.Name) != "" {
		if len(r.Name) > 255 {
			err.Append("name", lang.TT(l, "Maximo %v caracteres", 255))
		}
	} else {
		err.Append("name", lang.TT(l, "Este campo es requerido"))
	}
	return err.Errors()
}
