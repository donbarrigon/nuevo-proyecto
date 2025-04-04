package model

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/core"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Token struct {
	ID        bson.ObjectID `bson:"_id"`
	UserID    bson.ObjectID `bson:"user_id"`
	Token     string        `bson:"token"`
	CreatedAt time.Time     `bson:"created_at"`
	ExpiresAt time.Time     `bson:"expires_at"`
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

func (t *Token) Validate(ctx *core.Context) core.Error {
	if t.UserID.IsZero() {
		return &core.Err{
			Message: ctx.TT("Usuario invalido"),
			Err:     ctx.TT("El id de usuario esta vacio"),
		}
	}
	return nil
}
