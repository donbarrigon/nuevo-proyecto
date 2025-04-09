package services

import "github.com/donbarrigon/nuevo-proyecto/internal/app/controller"

func PermissionFilter(ctx *controller.Context) {
	allowFilters := map[string][]string{"name": {"eq", "ne", "gt", "gte", "lt", "lte"}}
	ctx.GetQueryFilter(allowFilters)
}
