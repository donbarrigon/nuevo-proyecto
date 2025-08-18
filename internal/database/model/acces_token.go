package model

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AccessToken struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      bson.ObjectID `bson:"user_id"       json:"user_id"`
	Token       string        `bson:"token"         json:"token"`
	Permissions []string      `bson:"permissions"   json:"permissions"`
	CreatedAt   time.Time     `bson:"created_at"    json:"created_at"`
	ExpiresAt   time.Time     `bson:"expires_at"    json:"expires_at"`
	app.Odm
}

func (t *AccessToken) CollectionName() string { return "access_tokens" }
func (t *AccessToken) GetID() bson.ObjectID   { return t.ID }
func (t *AccessToken) SetID(id bson.ObjectID) { t.ID = id }

func (t *AccessToken) BeforeCreate() app.Error {
	t.CreatedAt = time.Now()
	t.ExpiresAt = time.Now().Add(100 * time.Hour)
	return nil
}

func (t *AccessToken) BeforeUpdate() app.Error {
	t.ExpiresAt = time.Now().Add(100 * time.Hour)
	return nil
}

func NewAccessToken(userID bson.ObjectID, permissions []string) (*AccessToken, app.Error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		app.Log.Warning("Fail to create access token: " + err.Error())
	}
	tk := hex.EncodeToString(bytes)
	token := &AccessToken{
		//ID:        bson.NewObjectID(),
		UserID:      userID,
		Token:       tk,
		Permissions: permissions,
		CreatedAt:   time.Now(),
	}
	token.Refresh()
	token.Odm.Model = token

	err := token.Create()

	return token, err
}

func (t *AccessToken) Refresh() {
	t.ExpiresAt = time.Now().Add(1 * time.Hour)
}

func (t *AccessToken) Can(permissionNames ...string) app.Error {
	var result struct {
		ExpiresAt time.Time `bson:"expires_at"`
	}
	err := app.DB.Collection(t.CollectionName()).FindOne(context.TODO(), bson.D{
		{Key: "_id", Value: t.ID},
		{Key: "permissions", Value: bson.D{{Key: "$in", Value: permissionNames}}},
	},
		options.FindOne().SetProjection(bson.M{"expires_at": 1}),
	).Decode(&result)
	if err != nil {
		return app.Errors.Forbidden(err)
	}
	if result.ExpiresAt.Before(time.Now()) {
		return app.Errors.Forbiddenf("access denied: token expired at :expires_at", app.Item{Key: "expires_at", Value: result.ExpiresAt})
	}
	return nil
}

func (t *AccessToken) Anonymous() *AccessToken {
	var id bson.ObjectID // zero value: "000000000000000000000000"
	var timeZero time.Time
	return &AccessToken{
		ID:        id,
		UserID:    id,
		Token:     id.Hex(),
		CreatedAt: timeZero,
		ExpiresAt: timeZero,
	}
}
