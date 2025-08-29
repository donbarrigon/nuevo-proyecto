package model

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/qb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role struct {
	ID            bson.ObjectID   `bson:"_id,omitempty"         json:"id"`
	Name          string          `bson:"name"                  json:"name"`
	PermissionIDs []bson.ObjectID `bson:"permission_ids"        json:"-"`
	Permissions   []*Permission   `bson:"permissions,omitempty" json:"permissions,omitempty"` // manyToMany
	app.Odm       `bson:"-" json:"-"`
}

func NewRole() *Role {
	role := &Role{}
	role.Odm.Model = role
	return role
}

func (r *Role) CollectionName() string { return "roles" }
func (r *Role) GetID() bson.ObjectID   { return r.ID }
func (r *Role) SetID(id bson.ObjectID) { r.ID = id }

func (r *Role) BeforeCreate() app.Error { return nil }

func (r *Role) BeforeUpdate() app.Error { return nil }

func (r *Role) WithPermissions() bson.D {
	return qb.ManyToMany("permissions", "permission_ids")
}
