package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Profile struct {
	ID              bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID          bson.ObjectID  `bson:"user_id" json:"user_id"`
	Avatar          string         `bson:"avatar,omitempty" json:"avatar,omitempty"`
	FullName        string         `bson:"full_name,omitempty" json:"full_name,omitempty"`
	Nickname        string         `bson:"nickname" json:"nickname"`
	PhoneNumber     string         `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	DiscordUsername string         `bson:"discord_username,omitempty" json:"discord_username,omitempty"`
	CityID          string         `bson:"city_id" json:"city_id"`
	Preferences     map[string]any `bson:"preferences,omitempty" json:"preferences,omitempty"`
	CreatedAt       time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `bson:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time     `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
