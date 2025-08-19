package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type VerificationCode struct {
	ID        bson.ObjectID     `bson:"_id,omitempty"      json:"id,omitempty"`
	UserID    bson.ObjectID     `bson:"user_id"            json:"user_id"`
	Type      string            `bson:"type"               json:"type"`
	Code      string            `bson:"code"               json:"code"`
	Metadata  map[string]string `bson:"metadata,omitempty" json:"metadata,omitempty"`
	UsedAt    *time.Time        `bson:"used_at,omitempty"  json:"used_at,omitempty"`
	CreatedAt time.Time         `bson:"created_at"         json:"created_at"`
	ExpiresAt time.Time         `bson:"expires_at"         json:"updatexpires_ated_at"`

	app.Odm
}

func NewVerificationCode() *VerificationCode {
	verificationCode := &VerificationCode{}
	verificationCode.Odm.Model = verificationCode
	return verificationCode
}

func (v *VerificationCode) CollectionName() string { return "verification_codes" }
func (v *VerificationCode) GetID() bson.ObjectID   { return v.ID }
func (v *VerificationCode) SetID(id bson.ObjectID) { v.ID = id }

func (v *VerificationCode) BeforeCreate() app.Error {
	v.CreatedAt = time.Now()
	v.ExpiresAt = time.Now().Add(1 * time.Hour)
	return nil
}
func (v *VerificationCode) BeforeUpdate() app.Error { return nil }

func (v *VerificationCode) Generate(id bson.ObjectID, t string, metadata ...map[string]string) app.Error {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		app.Log.Warning("Fail to create verification code: " + err.Error())
		return app.Errors.InternalServerError(err)
	}
	code := hex.EncodeToString(bytes)

	if len(metadata) > 0 {
		v.Metadata = metadata[0]
	} else {
		v.Metadata = map[string]string{}
	}

	v.UserID = id
	v.Type = t
	v.Code = code

	return v.Create()
}
