package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	. "github.com/donbarrigon/nuevo-proyecto/internal/database/db/querybuilder"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID            bson.ObjectID   `bson:"_id,omitempty"         json:"id"`
	Email         string          `bson:"email"                 json:"email"`
	Password      string          `bson:"password"              json:"-"`
	AccessTokens  []*AccessToken  `bson:"tokens,omitempty"      json:"tokens,omitempty"`  // hasMany
	Profile       *Profile        `bson:"profile,omitempty"     json:"profile,omitempty"` //hasOne
	RoleIDs       []bson.ObjectID `bson:"role_ids"              json:"-"`
	Roles         []*Role         `bson:"roles,omitempty"       json:"roles,omitempty"` // manyToMany
	PermissionIDs []bson.ObjectID `bson:"permission_ids"        json:"-"`
	Permissions   []*Permission   `bson:"permissions,omitempty" json:"permissions,omitempty"` // manyToMany
	CreatedAt     time.Time       `bson:"created_at"            json:"created_at"`
	UpdatedAt     time.Time       `bson:"updated_at"            json:"updated_at"`
	DeletedAt     *time.Time      `bson:"deleted_at,omitempty"  json:"deleted_at,omitempty"`
}

func (u *User) CollectionName() string { return "users" }

func (u *User) GetID() bson.ObjectID { return u.ID }

func (u *User) SetID(id bson.ObjectID) { u.ID = id }

func (u *User) BeforeCreate() app.Error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate() app.Error {
	u.UpdatedAt = time.Now()
	return nil
}

// manyToMany
func (u *User) WithRoles() bson.D {
	return ManyToManyWith("roles", "role_ids", ManyToMany("permissions", "permission_ids"))
}

// manyToMany
func (u *User) WhithPermissions() bson.D {
	return ManyToMany("permissions", "permission_ids")
}

// hasMany
func (u *User) WithTokens() bson.D {
	return HasMany("tokens", "user_id")
}

// hasOne
func (u *User) WithProfile() []bson.D {
	return HasOne("profiles", "user_id", "profile")
}

func (u *User) Anonymous() *User {
	var id bson.ObjectID // zero value: "000000000000000000000000"
	var timeZero time.Time

	return &User{
		ID:        id,
		Email:     "anonymous@anonymous.com",
		Password:  "anonymous",
		CreatedAt: timeZero,
		UpdatedAt: timeZero,
		DeletedAt: &timeZero,
	}
}
