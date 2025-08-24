package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type City struct {
	ID          bson.ObjectID `bson:"_id,omitempty"          json:"id,omitempty"`
	Name        string        `bson:"name,omitempty"         json:"name,omitempty"`
	StateID     bson.ObjectID `bson:"state_id,omitempty"     json:"state_id,omitempty"`
	StateCode   string        `bson:"state_code,omitempty"   json:"state_code,omitempty"`
	StateName   string        `bson:"state_name,omitempty"   json:"state_name,omitempty"`
	CountryID   bson.ObjectID `bson:"country_id,omitempty"   json:"country_id,omitempty"`
	CountryCode string        `bson:"country_code,omitempty" json:"country_code,omitempty"`
	CountryName string        `bson:"country_name,omitempty" json:"country_name,omitempty"`
	Location    app.GeoPoint  `bson:"location"               json:"location"`
	Timezone    string        `bson:"timezone,omitempty"     json:"timezone,omitempty"`
	WikiDataID  string        `bson:"wikiDataId,omitempty"   json:"wikiDataId,omitempty"`
	CreatedAt   time.Time     `bson:"created_at"             json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"             json:"updated_at"`
	app.Odm     `bson:"-" json:"-"`
}

func NewCity() *City {
	city := &City{}
	city.Odm.Model = city
	return city
}

func (c *City) CollectionName() string { return "cities" }
func (c *City) GetID() bson.ObjectID   { return c.ID }
func (c *City) SetID(id bson.ObjectID) { c.ID = id }

func (c *City) BeforeCreate() app.Error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *City) BeforeUpdate() app.Error {
	c.UpdatedAt = time.Now()
	return nil
}
