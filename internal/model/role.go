package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Role struct {
	ID          bson.ObjectID `bson:"_id"`
	Name        string        `bson:"name"`
	Permissions []string      `bson:"permissions"`
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
func (r *Role) Index() map[string]string {
	return map[string]string{
		"name": "unique",
	}
}

//----------------------------------------------------------------//
type Permission struct {
	ID   bson.ObjectID `bson:"_id"`
	Name string        `bson:"name"`
}

func NewPermission() *Permission {
	return &Permission{
		ID: bson.NewObjectID(),
	}
}
func (p *Permission) CollectionName() string {
	return "roles"
}
func (p *Permission) GetID() bson.ObjectID {
	return p.ID
}
func (p *Permission) Index() map[string]string {
	return map[string]string{
		"name": "unique",
	}
}
