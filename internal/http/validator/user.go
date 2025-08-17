package validator

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/model"
)

type StoreUser struct {
	Name                 string         `json:"name"                       rules:"required|max:255"`
	Email                string         `json:"email"                      rules:"required_without:phone|email|max:255"`
	Password             string         `json:"password"                   rules:"reqired|confirmed|min:8|max:32"`
	PasswordConfirmation string         `json:"password_confirmation"`
	FullName             string         `json:"full_name,omitempty"        rules:"max:255"`
	Nickname             string         `json:"nickname"                   rules:"required|username|max:255"`
	PhoneNumber          string         `json:"phone_number,omitempty"     rules:"max:255"`
	DiscordUsername      string         `json:"discord_username,omitempty" rules:"max:255,regex:^(?!.*\\.\\.)(?!.*\\.$)(?!^\\.)[a-z0-9._]{2,32}$"`
	CityID               string         `json:"city_id"                    rules:"required|min:1"`
	Preferences          map[string]any `json:"preferences,omitempty"`
}

func (v *StoreUser) PrepareForValidation(ctx *app.HttpContext) app.Error {
	err := app.Errors.NewEmpty()
	city := &model.City{}
	if e := db.FindByHexID(city, v.CityID); e != nil {
		err.Appendf("city_id", "the city does not exist "+e.Error())
	} else {
		if city.ID.IsZero() {
			err.Appendf("city_id", "the city does not exist")
		}
	}
	return err
}

type UpdateUserEmail struct {
	Email string `json:"email" rules:"required_without:phone|email|max:255"`
}

type UpdateUserPassword struct {
	Password                string `json:"password" rules:"reqired"`
	NewPassword             string `json:"new_password" rules:"reqired|confirmed|min:8|max:32"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

type UpdateUserProfile struct {
	FullName        string         `json:"full_name,omitempty"        rules:"max:255"`
	Nickname        string         `json:"nickname"                   rules:"required|username|max:255"`
	PhoneNumber     string         `json:"phone_number,omitempty"     rules:"max:255"`
	DiscordUsername string         `json:"discord_username,omitempty" rules:"max:255,regex:^(?!.*\\.\\.)(?!.*\\.$)(?!^\\.)[a-z0-9._]{2,32}$"`
	CityID          string         `json:"city_id"                    rules:"required|min:1"`
	Preferences     map[string]any `json:"preferences,omitempty"`
}
