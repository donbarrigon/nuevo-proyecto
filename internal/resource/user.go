package resource

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
)

type UserLogin struct {
	AccessToken *model.AccessToken `json:"access_token"`
	User        *model.User        `json:"user"`
}

func NewUserLogin(u *model.User, t *model.AccessToken) *UserLogin {
	return &UserLogin{
		User:        u,
		AccessToken: t,
	}
}
