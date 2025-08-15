package controller

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/request"
)

func PermissionIndex(ctx *app.HttpContext) {
	allowFilters := map[string][]string{"name": {"eq", "ne", "lt", "gt", "lte", "gte", "sortable"}}

	var permissions []model.Permission
	res, err := db.Paginate(&model.Permission{}, &permissions, ctx.GetQueryFilter(allowFilters))
	if err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, res)
}

func PermissionExport(ctx *app.HttpContext) {
	allowFilters := map[string][]string{"name": {"eq", "ne", "lt", "gt", "lte", "gte", "sortable"}}
	qf := ctx.GetQueryFilter(allowFilters)
	qf.All()

	var permissions []model.Permission
	if err := db.FindByPipeline(&model.Permission{}, &permissions, qf.Pipeline()); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteCSV("db", permissions)
}

func PermissionShow(ctx *app.HttpContext) {
	id := ctx.Get("id")

	permission := &model.Permission{}
	if err := db.FindByHexID(permission, id); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, permission)
}

func PermissionStore(ctx *app.HttpContext) {
	req := &request.StorePermission{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.WriteError(err)
		return
	}

	permission := &model.Permission{}
	model.Fill(permission, req)

	if err := db.Create(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteCreated(permission)
}

func PermissionUpdate(ctx *app.HttpContext) {
	req := &request.StorePermission{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.WriteError(err)
		return
	}

	id := ctx.Get("id")
	permission := &model.Permission{}
	if err := db.FindByHexID(permission, id); err != nil {
		ctx.WriteError(err)
		return
	}

	model.Fill(permission, req)

	if err := db.Update(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteUpdated(permission)
}

func PermissionDestroy(ctx *app.HttpContext) {
	id := ctx.Get("id")
	permission := &model.Permission{}
	if err := db.FindByHexID(permission, id); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := db.Delete(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteDeleted(permission)
}

func PermissionRestore(ctx *app.HttpContext) {
	id := ctx.Get("id")
	permission := &model.Permission{}
	if err := db.FindByHexID(permission, id); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := db.Restore(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteRestored(permission)
}

func PermissionForceDelete(ctx *app.HttpContext) {
	id := ctx.Params["id"]
	permission := &model.Permission{}
	if err := db.FindByHexID(permission, id); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := db.ForceDelete(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteForceDeleted(permission)
}
