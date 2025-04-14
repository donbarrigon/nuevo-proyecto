package services

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/controller"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/resource"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
)

func PermissionFilter(ctx *controller.Context) {
	allowFilters := map[string][]string{
		"name":       {"eq", "ne", "lk", "gt", "gte", "lt", "lte", "sortable", "groupable"},
		"created_at": {"sortable"},
		"updated_at": {"sortable"},
	}
	qf := ctx.GetQueryFilter(allowFilters)

	var result []resource.Permission
	if err := db.Find(&model.Permission{}, &result, *qf.BsonD()); err != nil {
		ctx.WriteError(err)
	}

	ctx.WriteJSON(http.StatusOK, result)
}
