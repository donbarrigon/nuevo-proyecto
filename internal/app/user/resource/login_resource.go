package resource

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/user/model"
)

type TokenResource struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewTokenResource(t *model.Token) *TokenResource {
	return &TokenResource{
		ID:        t.ID.Hex(),
		UserID:    t.UserID.Hex(),
		Token:     t.Token,
		CreatedAt: t.CreatedAt,
		ExpiresAt: t.ExpiresAt,
	}
}

type UserLoginResource struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email,omitempty"`
	Phone     string         `json:"phone,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt time.Time      `json:"deletedAt,omitempty"`
	Token     *TokenResource `json:"token,omitempty"`
}

func NewUserLoginResource(u *model.User, t *model.Token) *UserLoginResource {
	var token *TokenResource
	// if t != nil {
	// 	token = NewTokenResource(t)
	// }
	return &UserLoginResource{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
		Token:     token,
	}
}
