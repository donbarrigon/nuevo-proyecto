package controller

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
)

func PermissionIndex(ctx *Context) {
	var permissions []*model.Permission

	if err := db.FindAll(&model.Permission{}, permissions); err != nil {
		ctx.WriteError(err)
		return
	}

	ctx.WriteJSON(http.StatusOK, permissions)
}
