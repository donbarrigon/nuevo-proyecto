package db

import (
	"context"

	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
)

type PaginateResource struct {
	// Campos comunes
	Data    any    `json:"data"`
	Path    string `json:"path"`
	PerPage int    `json:"per_page"`
	// Enlaces - estructura flexible que funciona para ambos tipos
	Links map[string]*string `json:"links"`

	// Campos para paginación por offset (nulos en paginación por cursor)
	CurrentPage *int `json:"current_page,omitempty"`
	LastPage    *int `json:"last_page,omitempty"`
	Total       *int `json:"total,omitempty"`
	From        *int `json:"from,omitempty"`
	To          *int `json:"to,omitempty"`

	// Campos para paginación por cursor (nulos en paginación por offset)
	NextCursor *string `json:"next_cursor,omitempty"`
	PrevCursor *string `json:"prev_cursor,omitempty"`
}

func Paginate(model MongoModel, result any, qf *QueryFilter) (*PaginateResource, errors.Error) {
	cursor, err := Mongo.Database.Collection(model.CollectionName()).Aggregate(context.TODO(), qf.Pipeline())
	if err != nil {
		return nil, errors.Mongo(err)
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return nil, errors.Mongo(err)
	}
	if qf.Page > 0 {
		// paginacion por off set
	} else {
		// paginacion por curscursor
	}

	return nil, nil
}
