package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type User struct {
	ID            bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name          string          `bson:"name" json:"name"`
	Email         string          `bson:"email" json:"email"`
	Password      string          `bson:"password" json:"-"`
	Tokens        []*Token        `bson:"tokens,omitempty" json:"tokens,omitempty"`   // hasMany
	Profile       *Profile        `bson:"profile,omitempty" json:"profile,omitempty"` //hasOne
	RoleIDs       []bson.ObjectID `bson:"role_ids" json:"-"`
	Roles         []*Role         `bson:"roles,omitempty" json:"roles,omitempty"` // manyToMany
	PermissionIDs []bson.ObjectID `bson:"permission_ids" json:"-"`
	Permissions   []*Permission   `bson:"permissions,omitempty" json:"permissions,omitempty"` // manyToMany
	CreatedAt     time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time       `bson:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time      `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func (u *User) CollectionName() string {
	return "users"
}

func (u *User) BefereCreate() app.Error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BefereUpdate() app.Error {
	u.UpdatedAt = time.Now()
	return nil
}
func (u *User) GetID() bson.ObjectID {
	return u.ID
}

func (u *User) SetID(id bson.ObjectID) {
	u.ID = id
}

// manyToMany
func (u *User) WithRoles() bson.D {
	return db.ManyToMany("roles", "role_ids")
}

// manyToMany
func (u *User) WhithPermissions() bson.D {
	return db.HasMany("permissions", "permission_ids")
}

// hasMany
func (u *User) WithTokens() bson.D {
	return db.HasMany("tokens", "user_id")
}

// hasOne
func (u *User) WithProfile() []bson.D {
	return db.HasOne("profiles", "user_id", "profile")
}

// manyToMany con preload de permissions
func (u *User) WithRolesAndPermissions() bson.D {
	return bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "roles"},
				{Key: "let", Value: bson.D{{Key: "role_ids", Value: "$role_ids"}}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$in", Value: bson.A{"$_id", "$$role_ids"}},
						}},
					}}},
					{{Key: "$lookup", Value: bson.D{
						{Key: "from", Value: "permissions"},
						{Key: "localField", Value: "permission_ids"},
						{Key: "foreignField", Value: "_id"},
						{Key: "as", Value: "permissions"},
					}}},
				}},
				{Key: "as", Value: "roles"},
			},
		},
	}
}

func (u *User) Anonymous() *User {
	var id bson.ObjectID // zero value: "000000000000000000000000"
	var timeZero time.Time

	return &User{
		ID:        id,
		Name:      "Anonymous",
		Email:     "anonymous@anonymous.com",
		Password:  "anonymous",
		CreatedAt: timeZero,
		UpdatedAt: timeZero,
		DeletedAt: &timeZero,
	}
}
