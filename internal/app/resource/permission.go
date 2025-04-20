package resource

import (
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Permission struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func NewPermission(permission *model.Permission) *Permission {
	return &Permission{
		ID:        permission.ID,
		Name:      permission.Name,
		DeletedAt: permission.DeletedAt,
	}
}
