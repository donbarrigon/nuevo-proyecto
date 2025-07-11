package model

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Token struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    bson.ObjectID `bson:"user_id"`
	Token     string        `bson:"token"`
	CreatedAt time.Time     `bson:"created_at"`
	ExpiresAt time.Time     `bson:"expires_at"`
}

func NewToken(userID bson.ObjectID) *Token {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("No se creo el token: " + err.Error())
	}
	token := hex.EncodeToString(bytes)
	tokenModel := &Token{
		//ID:        bson.NewObjectID(),
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
	}
	return tokenModel
}

func (t *Token) TableName() string {
	return "tokens"
}

func (t *Token) Default() {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	t.ExpiresAt = time.Now().Add(10 * time.Hour)
}

func (t *Token) Anonymous() *Token {
	var id bson.ObjectID // zero value: "000000000000000000000000"
	var timeZero time.Time
	return &Token{
		ID:        id,
		UserID:    id,
		Token:     id.Hex(),
		CreatedAt: timeZero,
		ExpiresAt: timeZero,
	}
}
