package controller

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	. "github.com/donbarrigon/nuevo-proyecto/internal/app/qb"
	"github.com/donbarrigon/nuevo-proyecto/internal/server/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/server/policy"
	"github.com/donbarrigon/nuevo-proyecto/internal/server/validator"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func RoleIndex(ctx *app.HttpContext) {
	if err := policy.RoleViewAny(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	role := model.NewRole()
	roles := []model.Role{}
	role.Aggregate(&roles, Pipeline(
		Match(),
		role.WithPermissions(),
	))

	ctx.ResponseOk(roles)
}

func RoleExport(ctx *app.HttpContext) {
	if err := policy.RoleViewAny(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	role := model.NewRole()
	roles := []model.Role{}
	role.Aggregate(&roles, Pipeline(
		Match(),
		role.WithPermissions(),
	))

	ctx.ResponseCSV("roles", roles)
}

func RoleTrashed(ctx *app.HttpContext) {
	if err := policy.RoleDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	trash := model.NewTrash()
	result := []model.Trash{}
	if err := trash.Find(&result, Filter(Where("collection", Eq("roles")))); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseOk(result)
}

func RoleShow(ctx *app.HttpContext) {
	if err := policy.RoleView(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	role := model.NewRole()
	if err := role.AggregateOne(Pipeline(
		Match(Where("_id", Eq(id))),
		role.WithPermissions(),
	)); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseOk(role)
}

func RoleStore(ctx *app.HttpContext) {
	if err := policy.RoleCreate(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	req := &validator.StoreRole{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	role := model.NewRole()
	if err := role.CreateBy(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), role, model.ACTION_CREATE, nil)

	ctx.ResponseCreated(role)
}

func RoleUpdate(ctx *app.HttpContext) {
	if err := policy.RoleUpdate(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	req := &validator.StoreRole{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	role := model.NewRole()
	if err := role.FindByHexID(ctx.Params["id"]); err != nil {
		ctx.ResponseError(err)
		return
	}

	original, _, err := role.UpdateBy(req)
	if err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), role, model.ACTION_UPDATE, original)

	ctx.ResponseOk(role)
}

func RoleDestroy(ctx *app.HttpContext) {
	if err := policy.RoleDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	role := model.NewRole()
	if err := role.FindByHexID(ctx.Params["id"]); err != nil {
		ctx.ResponseError(err)
		return
	}

	trash := model.NewTrash()
	if err := trash.MoveToTrash(role); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), role, model.ACTION_MOVE_TO_TRASH, nil)

	ctx.ResponseNoContent()
}

func RoleRestore(ctx *app.HttpContext) {
	if err := policy.RoleDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	role := model.NewRole()
	trahs := model.NewTrash()

	if err := trahs.RestoreByHexID(role, ctx.Params["id"]); err != nil {
		ctx.ResponseError(err)
		return
	}
	go model.HistoryRecord(ctx.Auth.GetUserID(), role, model.ACTION_RESTORE, nil)

	ctx.ResponseOk(role)
}

func RoleGrant(ctx *app.HttpContext) {
	req := &validator.GrantRole{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	role := model.NewRole()
	if err := role.First("_id", id); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := policy.RoleGrant(ctx, role); err != nil {
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

	user.RoleIDs = append(user.RoleIDs, role.ID)
	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), user, "grant", role)

	ctx.ResponseNoContent()
}

func RoleRevoke(ctx *app.HttpContext) {
	req := &validator.GrantRole{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	role := model.NewRole()
	if err := role.First("_id", id); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := policy.RoleRevoke(ctx, role); err != nil {
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

	for i, roleID := range user.RoleIDs {
		if roleID == role.ID {
			user.RoleIDs = append(user.RoleIDs[:i], user.RoleIDs[i+1:]...)
		}
	}

	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.HistoryRecord(ctx.Auth.GetUserID(), user, "revoke", role)

	ctx.ResponseNoContent()
}
