package resource

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/database/model"
)

type UserLogin struct {
	User        *model.User        `json:"user"`
	AccessToken *model.AccessToken `json:"access_token"`
}

func NewUserLogin(u *model.User, t *model.AccessToken) *UserLogin {
	return &UserLogin{
		User:        u,
		AccessToken: t,
	}
}
