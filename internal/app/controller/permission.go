package controller

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/app/resource"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
)

func PermissionIndex(ctx *Context) {
	allowFilters := map[string][]string{"name": {"eq", "ne", "lt", "gt", "lte", "gte", "sortable"}}

	var permissions []resource.Permission
	res, err := db.Paginate(&model.Permission{}, &permissions, ctx.GetQueryFilter(allowFilters))
	if err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, res)
}

func PermissionExport(ctx *Context) {
	allowFilters := map[string][]string{"name": {"eq", "ne", "lt", "gt", "lte", "gte", "sortable"}}
	qf := ctx.GetQueryFilter(allowFilters)
	qf.All()

	var permissions []resource.Permission
	if err := db.FindByPipeline(&model.Permission{}, permissions, qf.Pipeline()); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteCSV("bd", permissions)
}

func PremissionShow(ctx *Context) {
	id := ctx.LastParam()

	permission := &model.Permission{}
	if err := db.FindByHexID(permission, id); err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewPermission(permission)
	ctx.WriteJSON(http.StatusOK, res)
}

func PermissionStore(ctx *Context) {
	permission := &model.Permission{}
	if err := ctx.GetBody(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := permission.Validate(ctx.Lang()); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := db.Create(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewPermission(permission)
	ctx.WriteJSON(http.StatusOK, res)
}

func PermissionUpdate(ctx *Context) {
	req := &model.Permission{}
	if err := ctx.GetBody(req); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := req.Validate(ctx.Lang()); err != nil {
		ctx.WriteError(err)
		return
	}

	id := ctx.LastParam()
	permission := &model.Permission{}
	if err := db.FindByHexID(permission, id); err != nil {
		ctx.WriteError(err)
		return
	}

	Fill(permission, req)

	if err := db.Update(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewPermission(permission)
	ctx.WriteJSON(http.StatusOK, res)
}

func PermissionDelete(ctx *Context) {
	id := ctx.LastParam()
	permission := &model.Permission{}
	if err := db.FindByHexID(permission, id); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := db.Delete(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteNoContent()
}

func PermissionRestore(ctx *Context) {
	id := ctx.LastParam()
	permission := &model.Permission{}
	if err := db.FindByHexID(permission, id); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := db.Restore(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	res := resource.NewPermission(permission)
	ctx.WriteJSON(http.StatusOK, res)
}

func ForceDelete(ctx *Context) {
	id := ctx.LastParam()
	permission := &model.Permission{}
	if err := db.FindByHexID(permission, id); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := db.ForceDelete(permission); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteNoContent()
}
