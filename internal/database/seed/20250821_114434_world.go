package seed

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/shared/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func World() {
	app.PrintInfo("Seeding countries...")
	countriesIDs := map[int]bson.ObjectID{}
	statesIDd := map[int]bson.ObjectID{}
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
	app.PrintInfo("Seeded countries :total", app.E("total", totalCountries))

	app.PrintInfo("Seeding states...")
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
		statesIDd[seedState.ID] = state.GetID()
		totalStates++
	}
	for _, seedState := range seedStates {
		if seedState.ParentID != "" {
			state := model.NewState()
			state.FindByID(statesIDd[seedState.ID])
			var parentID int
			var er error
			parentID, er = strconv.Atoi(seedState.ParentID)
			if er != nil {
				app.PrintError("Fail to add parent state :error", app.E("error", er.Error()))
				panic(er.Error())
			}
			state.ParentID = statesIDd[parentID]
			if err := state.Update(); err != nil {
				app.PrintError("Fail to add parent state :error", app.E("error", err.Error()))
				panic(err.Error())
			}
		}
	}

	app.PrintInfo("Seeded states :total", app.E("total", totalStates))

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

func loadStates() []RegionSeed {
	filePath := "internal/database/json/states.json"
	data, er := os.ReadFile(filePath)
	if er != nil {
		app.PrintError("Fail to open :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	var states []RegionSeed
	if er := json.Unmarshal(data, &states); er != nil {
		app.PrintError("Fail to unmarshal :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	return states

}

func loadCities() []model.City {
	filePath := "internal/database/json/cities.json"
	data, er := os.ReadFile(filePath)
	if er != nil {
		app.PrintError("Fail to open :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	var cities []model.City
	if er := json.Unmarshal(data, &cities); er != nil {
		app.PrintError("Fail to unmarshal :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	return cities
}

// types

type CountrySeed struct {
	ID             int                     `json:"id"`
	Name           string                  `json:"name"`
	Iso3           string                  `json:"iso3"`
	Iso2           string                  `json:"iso2"`
	NumericCode    string                  `json:"numeric_code"`
	Phonecode      string                  `json:"phonecode"`
	Capital        string                  `json:"capital"`
	Currency       string                  `json:"currency"`
	CurrencyName   string                  `json:"currency_name"`
	CurrencySymbol string                  `json:"currency_symbol"`
	TLD            string                  `json:"tld"`
	Native         string                  `json:"native"`
	Region         string                  `json:"region"`
	RegionID       int                     `json:"region_id"`
	Subregion      string                  `json:"subregion"`
	SubregionID    int                     `json:"subregion_id"`
	Nationality    string                  `json:"nationality"`
	Timezones      []model.CountryTimezone `json:"timezones"`
	Translations   map[string]string       `json:"translations"`
	Latitude       string                  `json:"latitude"`
	Longitude      string                  `json:"longitude"`
	Emoji          string                  `json:"emoji"`
	EmojiU         string                  `json:"emojiU"`
}

type RegionSeed struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CountryID   int    `json:"country_id"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	Iso2        string `json:"iso2"`
	Iso3166_2   string `json:"iso3166_2"`
	FipsCode    string `json:"fips_code"`
	Type        string `json:"type"`
	Level       string `json:"level"`     // puede ser null
	ParentID    string `json:"parent_id"` // puede ser null
	Latitude    string `json:"latitude"`  // viene como string en JSON
	Longitude   string `json:"longitude"` // viene como string en JSON
	Timezone    string `json:"timezone"`
}
