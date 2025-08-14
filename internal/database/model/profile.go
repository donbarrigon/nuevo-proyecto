package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Profile struct {
	ID              bson.ObjectID  `bson:"_id,omitempty"              json:"id"`
	UserID          bson.ObjectID  `bson:"user_id"                    json:"user_id"`
	Avatar          string         `bson:"avatar,omitempty"           json:"avatar,omitempty"`
	FullName        string         `bson:"full_name,omitempty"        json:"full_name,omitempty"`
	Nickname        string         `bson:"nickname"                   json:"nickname"`
	PhoneNumber     string         `bson:"phone_number,omitempty"     json:"phone_number,omitempty"`
	DiscordUsername string         `bson:"discord_username,omitempty" json:"discord_username,omitempty"`
	CityID          string         `bson:"city_id"                    json:"city_id"`
	Preferences     map[string]any `bson:"preferences,omitempty"      json:"preferences,omitempty"`
	CreatedAt       time.Time      `bson:"created_at"                 json:"created_at"`
	UpdatedAt       time.Time      `bson:"updated_at"                 json:"updated_at"`
	DeletedAt       *time.Time     `bson:"deleted_at,omitempty"       json:"deleted_at,omitempty"`
}

func (p *Profile) CollectionName() string { return "profiles" }

func (p *Profile) GetID() bson.ObjectID { return p.ID }

func (p *Profile) SetID(id bson.ObjectID) { p.ID = id }

func (p *Profile) BeforeCreate() app.Error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Profile) BeforeUpdate() app.Error {
	p.UpdatedAt = time.Now()
	return nil
}
