package model

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Token struct {
	ID        bson.ObjectID `bson:"_id"`
	UserID    bson.ObjectID `bson:"userId"`
	Token     string        `bson:"token"`
	CreatedAt time.Time     `bson:"createdAt"`
	ExpiresAt time.Time     `bson:"expiresAt"`
}

func NewToken() *Token {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("No se creo el token: " + err.Error())
	}
	token := hex.EncodeToString(bytes)
	tokenModel := &Token{
		ID:        bson.NewObjectID(),
		Token:     token,
		CreatedAt: time.Now(),
	}
	tokenModel.Refresh()
	return tokenModel
}

func (t *Token) CollectionName() string {
	return "tokens"
}

func (t *Token) Refresh() {
	t.ExpiresAt = time.Now().Add(10 * time.Hour)
}

func (t *Token) GetID() bson.ObjectID {
	return t.ID
}

func (t *Token) Index() map[string]string {
	return map[string]string{
		"token":     "unique",
		"expiresAt": "index",
	}
}
