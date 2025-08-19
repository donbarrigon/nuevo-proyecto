package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type State struct {
	ID          bson.ObjectID `bson:"_id,omitempty"          json:"id,omitempty"`
	Name        string        `bson:"name,omitempty"         json:"name,omitempty"`
	CountryID   bson.ObjectID `bson:"country_id,omitempty"   json:"country_id,omitempty"`
	CountryCode string        `bson:"country_code,omitempty" json:"country_code,omitempty"`
	CountryName string        `bson:"country_name,omitempty" json:"country_name,omitempty"`
	ISO2        string        `bson:"iso2,omitempty"         json:"iso2,omitempty"`
	ISO3166_2   string        `bson:"iso3166_2,omitempty"    json:"iso3166_2,omitempty"`
	FIPSCode    string        `bson:"fips_code,omitempty"    json:"fips_code,omitempty"`
	Type        string        `bson:"type,omitempty"         json:"type,omitempty"`
	Level       *string       `bson:"level,omitempty"        json:"level,omitempty"`
	ParentID    *string       `bson:"parent_id,omitempty"    json:"parent_id,omitempty"`
	Latitude    string        `bson:"latitude,omitempty"     json:"latitude,omitempty"`
	Longitude   string        `bson:"longitude,omitempty"    json:"longitude,omitempty"`
	Timezone    string        `bson:"timezone,omitempty"     json:"timezone,omitempty"`
	CreatedAt   time.Time     `bson:"created_at"             json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"             json:"updated_at"`
	app.Odm
}

func NewState() *State {
	state := &State{}
	state.Odm.Model = state
	return state
}

func (s *State) CollectionName() string { return "states" }
func (s *State) GetID() bson.ObjectID   { return s.ID }
func (s *State) SetID(id bson.ObjectID) { s.ID = id }

func (s *State) BeforeCreate() app.Error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func (s *State) BeforeUpdate() app.Error {
	s.UpdatedAt = time.Now()
	return nil
}
