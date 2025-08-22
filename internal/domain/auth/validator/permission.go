package validator

import "github.com/donbarrigon/nuevo-proyecto/internal/app"

type StorePermission struct {
	Name string `json:"name" rules:"required|alpha_spaces|max:255"`
}

func (v *StorePermission) PrepareForValidation(ctx *app.HttpContext) app.Error { return nil }

type GrantPermission struct {
	UserID string `json:"user_id" rules:"required|exists:users,_id"`
}

func (v *GrantPermission) PrepareForValidation(ctx *app.HttpContext) app.Error { return nil }
