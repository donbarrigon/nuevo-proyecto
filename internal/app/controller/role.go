package controller

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func RoleIndex(ctx *Context) {
	var roles []*model.Role

	if err := db.FindAll(&model.Role{}, &roles); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, &roles)
}

func RoleShow(ctx *Context) {
	id := ctx.LastParam()

	role := &model.Role{}
	if err := db.FindByHexID(role, id); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, role)
}

func RoleStore(ctx *Context) {
	role := &model.Role{}
	if err := ctx.GetBody(role); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := role.Validate(ctx.Lang()); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := db.Create(role); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusCreated, role)
}

func RoleUpdate(ctx *Context) {
	id := ctx.LastParam()

	role := &model.Role{}
	if err := db.FindByHexID(role, id); err != nil {
		ctx.WriteError(err)
		return
	}

	req := &model.Role{}
	if err := ctx.GetBody(req); err != nil {
		ctx.WriteError(err)
		return
	}

	role.Name = req.Name

	if err := role.Validate(ctx.Lang()); err != nil {
		ctx.WriteError(err)
		return
	}

	if err := db.Update(role); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, role)
}

func RoleDestroy(ctx *Context) {
	id := ctx.LastParam()

	oid, e := bson.ObjectIDFromHex(id)
	if e != nil {
		ctx.WriteError(errors.HexID(e))
	}

	role := &model.Role{ID: oid}
	if err := db.Delete(role); err != nil {
		ctx.WriteError(err)
	}

	ctx.WriteNoContent()
}

func RoleRestore(ctx *Context) {
	id := ctx.LastParam()

	oid, e := bson.ObjectIDFromHex(id)
	if e != nil {
		ctx.WriteError(errors.HexID(e))
	}

	role := &model.Role{ID: oid}
	if err := db.Restore(role); err != nil {
		ctx.WriteError(err)
	}

	ctx.WriteNoContent()
}

func RoleForceDelete(ctx *Context) {

}

func RoleAppendPermission(ctx *Context) {
	id := ctx.LastParam()

	role := &model.Role{}
	if err := db.FindByHexID(role, id); err != nil {
		ctx.WriteError(err)
		return
	}

	req := make(map[string]string)
	if err := ctx.GetBody(req); err != nil {
		ctx.WriteError(err)
		return
	}

	permission := &model.Permission{}
	if err := db.FindByHexID(permission, req["permission_id"]); err != nil {
		ctx.WriteError(err)
		return
	}
	role.Permissions = append(role.Permissions, permission)

	if err := db.Update(role); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, role)
}

func RoleRemovePermission(ctx *Context) {
	id := ctx.LastParam()

	role := &model.Role{}
	if err := db.FindByHexID(role, id); err != nil {
		ctx.WriteError(err)
		return
	}

	req := make(map[string]string)
	if err := ctx.GetBody(req); err != nil {
		ctx.WriteError(err)
		return
	}

	for i, permission := range role.Permissions {
		if req["permission_id"] == permission.ID.Hex() {
			role.Permissions = append(role.Permissions[:i], role.Permissions[i+1:]...)
			break
		}
	}

	if err := db.Update(role); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, role)
}
