package controller

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	. "github.com/donbarrigon/nuevo-proyecto/internal/app/qb"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/policy"
	"github.com/donbarrigon/nuevo-proyecto/internal/validator"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func PermissionIndex(ctx *app.HttpContext) {
	if err := policy.PermissionViewAny(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	permission := model.NewPermission()
	result := []model.Permission{}
	if err := permission.Find(&result, GetAll(), FindOptions(ctx)); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseOk(result)
}

func PermissionExport(ctx *app.HttpContext) {
	if err := policy.PermissionViewAny(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	permission := model.NewPermission()
	result := []model.Permission{}
	if err := permission.Find(&result, GetAll()); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseCSV("permissions", result)
}

func PermissionTrashed(ctx *app.HttpContext) {
	if err := policy.PermissionDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	trash := model.NewTrash()
	result := []model.Trash{}
	if err := trash.Find(&result, Filter(Where("collection", Eq("permissions")))); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseOk(result)
}

func PermissionShow(ctx *app.HttpContext) {
	if err := policy.PermissionView(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	permission := model.NewPermission()
	if err := permission.FindByHexID(ctx.Params["id"]); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseOk(permission)
}

func PermissionStore(ctx *app.HttpContext) {
	if err := policy.PermissionCreate(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	req := &validator.StorePermission{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	permission := model.NewPermission()
	if err := permission.CreateBy(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), permission, model.ACTION_CREATE, nil)

	ctx.ResponseCreated(permission)
}

func PermissionUpdate(ctx *app.HttpContext) {
	if err := policy.PermissionUpdate(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	req := &validator.StorePermission{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	permission := model.NewPermission()
	if err := permission.FindByHexID(ctx.Params["id"]); err != nil {
		ctx.ResponseError(err)
		return
	}

	original, _, err := permission.UpdateBy(req)
	if err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), permission, model.ACTION_UPDATE, original)

	ctx.ResponseOk(permission)
}

func PermissionDestroy(ctx *app.HttpContext) {
	if err := policy.PermissionDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	permission := model.NewPermission()
	if err := permission.FindByHexID(ctx.Params["id"]); err != nil {
		ctx.ResponseError(err)
		return
	}

	trash := model.NewTrash()
	if err := trash.MoveToTrash(permission); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), permission, model.ACTION_DELETE, nil)

	ctx.ResponseNoContent()
}

func PermissionRestore(ctx *app.HttpContext) {
	if err := policy.PermissionDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	permission := model.NewPermission()
	trash := model.NewTrash()

	if err := trash.RestoreByHexID(permission, ctx.Params["id"]); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), permission, model.ACTION_RESTORE, nil)

	ctx.ResponseOk(permission)
}

func PermissionGrant(ctx *app.HttpContext) {
	req := &validator.GrantPermission{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	permission := model.NewPermission()
	if err := permission.First("_id", id); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := policy.PermissionGrant(ctx, permission); err != nil {
		ctx.ResponseError(err)
		return
	}

	userID, err := bson.ObjectIDFromHex(req.UserID)
	if err != nil {
		ctx.ResponseError(app.Errors.HexID(err))
		return
	}

	user := model.NewUser()
	if err := user.First("_id", userID); err != nil {
		ctx.ResponseError(err)
		return
	}

	user.PermissionIDs = append(user.PermissionIDs, permission.ID)
	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), user, "grant", permission)

	ctx.ResponseNoContent()
}

func PermissionRevoke(ctx *app.HttpContext) {
	req := &validator.GrantPermission{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	permission := model.NewPermission()
	if err := permission.First("_id", id); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := policy.PermissionRevoke(ctx, permission); err != nil {
		ctx.ResponseError(err)
		return
	}

	userID, err := bson.ObjectIDFromHex(req.UserID)
	if err != nil {
		ctx.ResponseError(app.Errors.HexID(err))
		return
	}

	user := model.NewUser()
	if err := user.First("_id", userID); err != nil {
		ctx.ResponseError(err)
		return
	}

	for i, p := range user.PermissionIDs {
		if p == permission.ID {
			user.PermissionIDs = append(user.PermissionIDs[:i], user.PermissionIDs[i+1:]...)
		}
	}

	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), user, "revoke", permission)

	ctx.ResponseNoContent()
}
