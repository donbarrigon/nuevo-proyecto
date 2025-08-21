package validator

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
)

type StoreUser struct {
	Email                string         `json:"email"                      rules:"required|max:255|email|unique:users,email"`
	Password             string         `json:"password"                   rules:"reqired|confirmed|min:8|max:32"`
	PasswordConfirmation string         `json:"password_confirmation"`
	FullName             string         `json:"full_name,omitempty"        rules:"max:255|alpha_spaces_accents"`
	Nickname             string         `json:"nickname"                   rules:"required|username|max:255"`
	PhoneNumber          string         `json:"phone_number,omitempty"     rules:"max:255|regex:^\\+[1-9]\\d{1,14}$"`
	DiscordUsername      string         `json:"discord_username,omitempty" rules:"max:255"`
	CityID               string         `json:"city_id"                    rules:"required|max:255|exists:cities,_id"`
	Preferences          map[string]any `json:"preferences,omitempty"`
}

func (v *StoreUser) PrepareForValidation(ctx *app.HttpContext) app.Error {
	err := app.Errors.NewEmpty()
	// city := model.NewCity()
	// if e := city.FindByHexID(v.CityID); e != nil {
	// 	err.Appendf("city_id", "the city does not exist")
	// } else {
	// 	if city.ID.IsZero() {
	// 		err.Appendf("city_id", "the city does not exist")
	// 	}
	// }
	return err
}

type UpdateUserEmail struct {
	Email    string `json:"email" rules:"required|max:255|email|unique:users,email"`
	Password string `json:"password" rules:"required"`
}

func (v *UpdateUserEmail) PrepareForValidation(ctx *app.HttpContext) app.Error { return nil }

type UpdateUserPassword struct {
	Password                string `json:"password" rules:"reqired"`
	NewPassword             string `json:"new_password" rules:"reqired|confirmed|min:8|max:32"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

func (v *UpdateUserPassword) PrepareForValidation(ctx *app.HttpContext) app.Error { return nil }

type UpdateUserProfile struct {
	FullName        string         `json:"full_name,omitempty"        rules:"max:255|alpha_spaces_accents"`
	Nickname        string         `json:"nickname"                   rules:"required|username|max:255"`
	PhoneNumber     string         `json:"phone_number,omitempty"     rules:"max:255"`
	DiscordUsername string         `json:"discord_username,omitempty" rules:"max:255,regex:^(?!.*\\.\\.)(?!.*\\.$)(?!^\\.)[a-z0-9._]{2,32}$"`
	CityID          string         `json:"city_id"                    rules:"required|min:1"`
	Preferences     map[string]any `json:"preferences,omitempty"`
}

func (v *UpdateUserProfile) PrepareForValidation(ctx *app.HttpContext) app.Error { return nil }

type UserLogin struct {
	Email    string `json:"email"    rules:"required|email"`
	Password string `json:"password" rules:"required"`
}

func (v *UserLogin) PrepareForValidation(ctx *app.HttpContext) app.Error { return nil }
