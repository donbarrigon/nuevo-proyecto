package model

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Country struct {
	ID             bson.ObjectID     `bson:"id,omitempty"              json:"id,omitempty"`
	Name           string            `bson:"name,omitempty"            json:"name,omitempty"`
	Iso3           string            `bson:"iso3,omitempty"            json:"iso3,omitempty"`
	Iso2           string            `bson:"iso2,omitempty"            json:"iso2,omitempty"`
	NumericCode    string            `bson:"numeric_code,omitempty"    json:"numeric_code,omitempty"`
	PhoneCode      string            `bson:"phonecode,omitempty"       json:"phonecode,omitempty"`
	Capital        string            `bson:"capital,omitempty"         json:"capital,omitempty"`
	Currency       string            `bson:"currency,omitempty"        json:"currency,omitempty"`
	CurrencyName   string            `bson:"currency_name,omitempty"   json:"currency_name,omitempty"`
	CurrencySymbol string            `bson:"currency_symbol,omitempty" json:"currency_symbol,omitempty"`
	TLD            string            `bson:"tld,omitempty"             json:"tld,omitempty"`
	Native         string            `bson:"native,omitempty"          json:"native,omitempty"`
	Region         CountryRegion     `bson:"region,omitempty"          json:"region,omitempty"`
	Subregion      CountrySubRegion  `bson:"subregion,omitempty"       json:"subregion,omitempty"`
	Nationality    string            `bson:"nationality,omitempty"     json:"nationality,omitempty"`
	Timezones      []CountryTimezone `bson:"timezones,omitempty"       json:"timezones,omitempty"`
	Translations   map[string]string `bson:"translations,omitempty"    json:"translations,omitempty"`
	Location       app.GeoPoint      `bson:"location"                  json:"location"`
	Emoji          string            `bson:"emoji,omitempty"           json:"emoji,omitempty"`
	EmojiU         string            `bson:"emojiU,omitempty"          json:"emojiU,omitempty"`
	CreatedAt      time.Time         `bson:"created_at"                json:"created_at"`
	UpdatedAt      time.Time         `bson:"updated_at"                json:"updated_at"`
	app.Odm        `bson:"-" json:"-"`
}

type CountryTimezone struct {
	ZoneName      string `bson:"zoneName,omitempty"      json:"zoneName,omitempty"`
	GMTOffset     int    `bson:"gmtOffset,omitempty"     json:"gmtOffset,omitempty"`
	GMTOffsetName string `bson:"gmtOffsetName,omitempty" json:"gmtOffsetName,omitempty"`
	Abbreviation  string `bson:"abbreviation,omitempty"  json:"abbreviation,omitempty"`
	TZName        string `bson:"tzName,omitempty"        json:"tzName,omitempty"`
}

type CountryRegion struct {
	ID           int               `bson:"id,omitempty"           json:"id,omitempty"`
	Name         string            `bson:"name,omitempty"         json:"name,omitempty"`
	Translations map[string]string `bson:"translations,omitempty" json:"translations,omitempty"`
	WikiDataId   string            `bson:"wikiDataId,omitempty"   json:"wikiDataId,omitempty"`
}

type CountrySubRegion struct {
	ID           int               `bson:"id,omitempty"           json:"id,omitempty"`
	RegionID     int               `bson:"region_id,omitempty"    json:"region_id,omitempty"`
	Name         string            `bson:"name,omitempty"         json:"name,omitempty"`
	Translations map[string]string `bson:"translations,omitempty" json:"translations,omitempty"`
	WikiDataId   string            `bson:"wikiDataId,omitempty"   json:"wikiDataId,omitempty"`
}

func NewCountry() *Country {
	country := &Country{}
	country.Odm.Model = country
	return country
}

func (c *Country) CollectionName() string { return "countries" }
func (c *Country) GetID() bson.ObjectID   { return c.ID }
func (c *Country) SetID(id bson.ObjectID) { c.ID = id }

func (c *Country) BeforeCreate() app.Error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Country) BeforeUpdate() app.Error {
	c.UpdatedAt = time.Now()
	return nil
}
