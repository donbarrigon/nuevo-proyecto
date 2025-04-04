package user

import "github.com/donbarrigon/nuevo-proyecto/internal/model"

type RoleResouce struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func NewRoleResouce(r *model.Role) *RoleResouce {

	return &RoleResouce{
		ID:          r.ID.Hex(),
		Name:        r.Name,
		Permissions: r.Permissions,
	}
}
