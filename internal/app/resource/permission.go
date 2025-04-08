package resource

import "github.com/donbarrigon/nuevo-proyecto/internal/app/model"

type Permission struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewPermission(permission *model.Permission) *Permission {
	return &Permission{
		ID:   permission.ID.Hex(),
		Name: permission.Name,
	}
}

func NewPermissionCollection(permissions []*model.Permission) []*Permission {
	result := make([]*Permission, 0)
	for _, model := range permissions {
		result = append(result, NewPermission(model))
	}
	return result
}
