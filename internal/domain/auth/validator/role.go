package validator

import "github.com/donbarrigon/nuevo-proyecto/internal/app"

type StoreRole struct {
	Name string `json:"name" rules:"required|alpha_spaces|max:255"`
}

func (v *StoreRole) PrepareForValidation(ctx *app.HttpContext) app.Error { return nil }

type GrantRole struct {
	UserID string `json:"user_id" rules:"required|exists:users,_id"`
}

func (v *GrantRole) PrepareForValidation(ctx *app.HttpContext) app.Error { return nil }
