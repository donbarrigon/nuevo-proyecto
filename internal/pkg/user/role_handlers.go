package user

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

func StoreRole(ctx *app.HandlerContext) {
	req := &RoleRequest{}
	if err := app.GetRequest(ctx, req, http.MethodPost); err != nil {
		err.WriteResponse(ctx.Writer)
		return
	}

	roleModel := &model.Role{Name: req.Name}
	if _, err := app.Mongo.Create(roleModel); err != nil {
		e := app.ErrorJSON{
			Message: lang.M(ctx.Lang(), "app.internal-error"),
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		}
		e.WriteResponse(ctx.Writer)
		return
	}

	app.ResponseJSON(ctx.Writer, NewRoleResouce(roleModel), 200)
}
