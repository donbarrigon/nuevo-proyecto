package model

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AccessToken struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"user_id"       json:"user_id"`
	Token     string        `bson:"token"         json:"token"`
	CreatedAt time.Time     `bson:"created_at"    json:"created_at"`
	ExpiresAt time.Time     `bson:"expires_at"    json:"expires_at"`
}

func (t *AccessToken) CollectionName() string { return "tokens" }

func (t *AccessToken) GetID() bson.ObjectID { return t.ID }

func (t *AccessToken) SetID(id bson.ObjectID) { t.ID = id }

func (t *AccessToken) BeforeCreate() app.Error {
	t.CreatedAt = time.Now()
	t.ExpiresAt = time.Now().Add(100 * time.Hour)
	return nil
}

func (t *AccessToken) BefereUpdate() app.Error {
	t.ExpiresAt = time.Now().Add(100 * time.Hour)
	return nil
}

func NewAccessToken(userID bson.ObjectID) *AccessToken {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("No se creo el token de: " + err.Error())
	}
	token := hex.EncodeToString(bytes)
	tokenModel := &AccessToken{
		//ID:        bson.NewObjectID(),
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
	}
	return tokenModel
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
