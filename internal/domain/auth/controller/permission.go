package controller

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	. "github.com/donbarrigon/nuevo-proyecto/internal/app/qb"
	"github.com/donbarrigon/nuevo-proyecto/internal/domain/auth/policy"
	"github.com/donbarrigon/nuevo-proyecto/internal/domain/auth/validator"
	"github.com/donbarrigon/nuevo-proyecto/internal/shared/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/shared/service"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func PermissionIndex(ctx *app.HttpContext) {
	if err := policy.PermissionViewAny(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	permission := model.NewPermission()
	permissions := []model.Permission{}
	if err := permission.Find(&permissions, Document()); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseOk(permissions)
}

func PermissionExport(ctx *app.HttpContext) {
	if err := policy.PermissionViewAny(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	permission := model.NewPermission()
	permissions := []model.Permission{}
	if err := permission.Find(&permissions, Document()); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseCSV("permissions", permissions)
}

func PermissionTrashed(ctx *app.HttpContext) {
	if err := policy.PermissionDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	activity := model.NewActivity()
	result := []*model.Activity{}
	if err := activity.Find(result, Document(
		Where("collection", Eq("permissions")),
		Where("action", Eq("delete")),
	)); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseOk(activity)

}

func PermissionShow(ctx *app.HttpContext) {
	if err := policy.PermissionView(ctx); err != nil {
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

	go service.ActivityRecord(ctx.Auth.GetUserID(), permission, "create", permission)

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

	dirty, err := permission.UpdateBy(req)
	if err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.ActivityRecord(ctx.Auth.GetUserID(), permission, "update", dirty)

	ctx.ResponseOk(permission)
}

func PermissionDestroy(ctx *app.HttpContext) {
	if err := policy.PermissionDelete(ctx); err != nil {
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

	if err := permission.Delete(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.ActivityRecord(ctx.Auth.GetUserID(), permission, "delete", permission)

	ctx.ResponseNoContent()
}

func PermissionRestore(ctx *app.HttpContext) {
	if err := policy.PermissionDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	permission := model.NewPermission()
	activity := model.NewActivity()

	if err := activity.FindOne(Document(
		Where("document_id", Eq(id)),
		Where("collection", Eq(permission.CollectionName())),
		Where("action", Eq("delete")),
	)); err != nil {
		ctx.ResponseError(err)
		return
	}

	changes, ok := activity.Changes.(map[string]any)
	if !ok {
		ctx.ResponseError(app.Errors.InternalServerErrorf("invalid activity changes"))
		return
	}

	permission.ID = activity.DocumentID
	if err := app.FillByMap(permission, changes); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := permission.Create(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.ActivityRecord(ctx.Auth.GetUserID(), permission, "restore", permission)

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

	go service.ActivityRecord(ctx.Auth.GetUserID(), user, "grant", permission)

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

	go service.ActivityRecord(ctx.Auth.GetUserID(), user, "revoke", permission)

	ctx.ResponseNoContent()
}
