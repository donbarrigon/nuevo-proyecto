package model

import (
	"context"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	. "github.com/donbarrigon/nuevo-proyecto/internal/database/db/querybuilder"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (u *User) Can(permissionName string) app.Error {
	usersCol := db.Mongo.Database.Collection(u.CollectionName())

	// Pipeline para buscar si el usuario tiene el permiso
	pipeline := mongo.Pipeline{
		// Filtramos por el ID del usuario
		{{Key: "$match", Value: bson.M{"_id": u.ID}}},

		// Unimos permisos directos
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         "permissions",
				"localField":   "permission_ids",
				"foreignField": "_id",
				"as":           "direct_permissions",
			},
		}},

		// Unimos roles
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         "roles",
				"localField":   "role_ids",
				"foreignField": "_id",
				"as":           "roles",
			},
		}},

		// Desanidamos roles para unir sus permisos
		{{Key: "$unwind", Value: bson.M{"path": "$roles", "preserveNullAndEmptyArrays": true}}},

		// Unimos permisos de roles
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         "permissions",
				"localField":   "roles.permission_ids",
				"foreignField": "_id",
				"as":           "role_permissions",
			},
		}},

		// Agrupamos todo de nuevo (porque hicimos unwind)
		{{
			Key: "$group",
			Value: bson.M{
				"_id":                "$_id",
				"direct_permissions": bson.M{"$first": "$direct_permissions"},
				"role_permissions":   bson.M{"$push": "$role_permissions"},
			},
		}},

		// Flatten de role_permissions (de array de arrays → array)
		{{
			Key: "$project",
			Value: bson.M{
				"permissions": bson.M{
					"$setUnion": []interface{}{
						"$direct_permissions",
						bson.M{"$reduce": bson.M{
							"input":        "$role_permissions",
							"initialValue": bson.A{},
							"in":           bson.M{"$setUnion": []interface{}{"$$value", "$$this"}},
						}},
					},
				},
			},
		}},

		// Filtramos para ver si existe el permiso buscado
		{{
			Key: "$project",
			Value: bson.M{
				"hasPermission": bson.M{
					"$gt": bson.A{
						bson.M{"$size": bson.M{
							"$filter": bson.M{
								"input": "$permissions",
								"as":    "perm",
								"cond":  bson.M{"$eq": bson.A{"$$perm.name", permissionName}},
							},
						}},
						0,
					},
				},
			},
		}},
	}

	ctx := context.TODO()
	cursor, err := usersCol.Aggregate(ctx, pipeline)
	if err != nil {
		return app.Errors.Mongo(err)
	}
	defer cursor.Close(ctx)

	var result struct {
		HasPermission bool `bson:"hasPermission"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return app.Errors.Mongo(err)
		}
		if result.HasPermission {
			return nil
		}
	}

	return app.Errors.Forbiddenf("access denied")
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
