package resource

import "github.com/donbarrigon/nuevo-proyecto/internal/app/model"

type Permission struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

func NewPermission(permission *model.Permission) *Permission {
	return &Permission{
		ID:   permission.ID.Hex(),
		Name: permission.Name,
	}
}
