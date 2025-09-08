package seed

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func World() {
	app.PrintInfo("Seeding countries...")
	countriesIDs := map[int]bson.ObjectID{}
	statesIDs := map[int]bson.ObjectID{}
	regions := loadRegions()
	subregions := loadSubRegions()
	seedCountries := loadCountries()

	var totalCountries int
	for _, seedCountry := range seedCountries {
		country := model.NewCountry()
		country.Name = seedCountry.Name
		country.Iso3 = seedCountry.Iso3
		country.Iso2 = seedCountry.Iso2
		country.NumericCode = seedCountry.NumericCode
		country.PhoneCode = seedCountry.Phonecode
		country.Capital = seedCountry.Capital
		country.Currency = seedCountry.Currency
		country.CurrencyName = seedCountry.CurrencyName
		country.CurrencySymbol = seedCountry.CurrencySymbol
		country.TLD = seedCountry.TLD
		country.Native = seedCountry.Native
		country.Nationality = seedCountry.Nationality
		country.Timezones = seedCountry.Timezones
		country.Translations = seedCountry.Translations
		country.Emoji = seedCountry.Emoji
		country.EmojiU = seedCountry.EmojiU

		for _, region := range regions {
			if region.ID == seedCountry.RegionID {
				country.Region = region
				break
			}
		}
		for _, subregion := range subregions {
			if subregion.ID == seedCountry.SubregionID {
				country.Subregion = subregion
				break
			}
		}
		var latitude, longitud float64
		var er error
		latitude, er = strconv.ParseFloat(seedCountry.Latitude, 64)
		if er != nil {
			latitude = 0
		}
		longitud, er = strconv.ParseFloat(seedCountry.Longitude, 64)
		if er != nil {
			longitud = 0
		}
		country.Location = *app.NewGeoPoint(latitude, longitud)

		if err := country.Create(); err != nil {
			app.PrintError("Fail to create country :error", app.E("error", err.Error()))
			panic(err.Error())
		}
		countriesIDs[seedCountry.ID] = country.GetID()
		totalCountries++
	}
	app.PrintInfo("Seeded countries :total wait seeding states", app.E("total", totalCountries))

	seedStates := loadStates()
	var totalStates int
	for _, seedState := range seedStates {
		state := model.NewState()
		state.Name = seedState.Name
		state.CountryID = countriesIDs[seedState.CountryID]
		state.CountryCode = seedState.CountryCode
		state.CountryName = seedState.CountryName
		state.Iso2 = seedState.Iso2
		state.Iso3166_2 = seedState.Iso3166_2
		state.FipsCode = seedState.FipsCode
		state.Type = seedState.Type
		//state.ParentID = seedState.ParentID
		state.Timezone = seedState.Timezone

		var level int
		if seedState.Level != "" {
			level, _ = strconv.Atoi(seedState.Level)
		}
		if level == 0 {
			level = 1
		}
		state.Level = level

		var latitude, longitud float64
		var er error
		latitude, er = strconv.ParseFloat(seedState.Latitude, 64)
		if er != nil {
			latitude = 0
		}
		longitud, er = strconv.ParseFloat(seedState.Longitude, 64)
		if er != nil {
			longitud = 0
		}
		state.Location = *app.NewGeoPoint(latitude, longitud)

		if err := state.Create(); err != nil {
			app.PrintError("Fail to create state :error", app.E("error", err.Error()))
			panic(err.Error())
		}
		statesIDs[seedState.ID] = state.GetID()
		totalStates++
	}
	for _, seedState := range seedStates {
		if seedState.ParentID != "" {
			state := model.NewState()
			state.FindByID(statesIDs[seedState.ID])
			var parentID int
			var er error
			parentID, er = strconv.Atoi(seedState.ParentID)
			if er != nil {
				app.PrintError("Fail to add parent state :error", app.E("error", er.Error()))
				panic(er.Error())
			}
			state.ParentID = statesIDs[parentID]
			if err := state.Update(); err != nil {
				app.PrintError("Fail to add parent state :error", app.E("error", err.Error()))
				panic(err.Error())
			}
		}
	}
	app.PrintInfo("Seeded states :total wait seeding cities", app.E("total", totalStates))

	seedCities := loadCities()
	allcityes := []*model.City{}
	for _, seedCity := range seedCities {
		city := model.NewCity()
		city.Name = seedCity.Name
		city.StateID = statesIDs[seedCity.StateID]
		city.StateCode = seedCity.StateCode
		city.StateName = seedCity.StateName
		city.CountryID = countriesIDs[seedCity.CountryID]
		city.CountryCode = seedCity.CountryCode
		city.CountryName = seedCity.CountryName
		city.Timezone = seedCity.Timezone

		var latitude, longitud float64
		var er error
		latitude, er = strconv.ParseFloat(seedCity.Latitude, 64)
		if er != nil {
			latitude = 0
		}
		longitud, er = strconv.ParseFloat(seedCity.Longitude, 64)
		if er != nil {
			longitud = 0
		}
		city.Location = *app.NewGeoPoint(latitude, longitud)
		allcityes = append(allcityes, city)
	}
	app.PrintInfo("file cities ready :total", app.E("total", len(allcityes)))
	city := model.NewCity()
	if err := city.CreateMany(allcityes); err != nil {
		app.PrintError("Fail to create cities :error", app.E("error", err.Error()))
		panic(err.Error())
	}
	app.PrintInfo("Finish seed cities :total", app.E("total", len(allcityes)))

}

func loadRegions() []model.CountryRegion {

	filePath := "internal/database/json/regions.json"
	data, er := os.ReadFile(filePath)
	if er != nil {
		app.PrintError("Fail to open :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	var regions []model.CountryRegion
	if er := json.Unmarshal(data, &regions); er != nil {
		app.PrintError("Fail to unmarshal :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	return regions
}

func loadSubRegions() []model.CountrySubRegion {
	filePath := "internal/database/json/subregions.json"
	data, er := os.ReadFile(filePath)
	if er != nil {
		app.PrintError("Fail to open :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	var subRegions []model.CountrySubRegion
	if er := json.Unmarshal(data, &subRegions); er != nil {
		app.PrintError("Fail to unmarshal :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	return subRegions
}

func loadCountries() []CountrySeed {
	filePath := "internal/database/json/countries.json"
	data, er := os.ReadFile(filePath)
	if er != nil {
		app.PrintError("Fail to open :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	var countries []CountrySeed
	if er := json.Unmarshal(data, &countries); er != nil {
		app.PrintError("Fail to unmarshal :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	return countries
}

func loadStates() []StateSeed {
	filePath := "internal/database/json/states.json"
	data, er := os.ReadFile(filePath)
	if er != nil {
		app.PrintError("Fail to open :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	var states []StateSeed
	if er := json.Unmarshal(data, &states); er != nil {
		app.PrintError("Fail to unmarshal :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	return states

}

func loadCities() []CitySeed {
	filePath := "internal/database/json/cities.json"
	data, er := os.ReadFile(filePath)
	if er != nil {
		app.PrintError("Fail to open :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	var cities []CitySeed
	if er := json.Unmarshal(data, &cities); er != nil {
		app.PrintError("Fail to unmarshal :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	return cities
}

// types

type CountrySeed struct {
	ID             int                     `json:"id,omitempty"`
	Name           string                  `json:"name,omitempty"`
	Iso3           string                  `json:"iso3,omitempty"`
	Iso2           string                  `json:"iso2,omitempty"`
	NumericCode    string                  `json:"numeric_code,omitempty"`
	Phonecode      string                  `json:"phonecode,omitempty"`
	Capital        string                  `json:"capital,omitempty"`
	Currency       string                  `json:"currency,omitempty"`
	CurrencyName   string                  `json:"currency_name,omitempty"`
	CurrencySymbol string                  `json:"currency_symbol,omitempty"`
	TLD            string                  `json:"tld,omitempty"`
	Native         string                  `json:"native,omitempty"`
	Region         string                  `json:"region,omitempty"`
	RegionID       int                     `json:"region_id,omitempty"`
	Subregion      string                  `json:"subregion,omitempty"`
	SubregionID    int                     `json:"subregion_id,omitempty"`
	Nationality    string                  `json:"nationality,omitempty"`
	Timezones      []model.CountryTimezone `json:"timezones,omitempty"`
	Translations   map[string]string       `json:"translations,omitempty"`
	Latitude       string                  `json:"latitude,omitempty"`
	Longitude      string                  `json:"longitude,omitempty"`
	Emoji          string                  `json:"emoji,omitempty"`
	EmojiU         string                  `json:"emojiU,omitempty"`
}

type StateSeed struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	CountryID   int    `json:"country_id,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	CountryName string `json:"country_name,omitempty"`
	Iso2        string `json:"iso2,omitempty"`
	Iso3166_2   string `json:"iso3166_2,omitempty"`
	FipsCode    string `json:"fips_code,omitempty"`
	Type        string `json:"type,omitempty"`
	Level       string `json:"level,omitempty"`     // puede ser null
	ParentID    string `json:"parent_id,omitempty"` // puede ser null
	Latitude    string `json:"latitude,omitempty"`  // viene como string en JSON
	Longitude   string `json:"longitude,omitempty"` // viene como string en JSON
	Timezone    string `json:"timezone,omitempty"`
}

type CitySeed struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	StateID     int    `json:"state_id,omitempty"`
	StateCode   string `json:"state_code,omitempty"`
	StateName   string `json:"state_name,omitempty"`
	CountryID   int    `json:"country_id,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	CountryName string `json:"country_name,omitempty"`
	Latitude    string `json:"latitude,omitempty"`
	Longitude   string `json:"longitude,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
	WikiDataID  string `json:"wikiDataId,omitempty"`
}
