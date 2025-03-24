package user

import (
	"net/http"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
	"github.com/donbarrigon/nuevo-proyecto/pkg/validate"
	"golang.org/x/crypto/bcrypt"
)

func storeService(ctx *app.ControllerContext, req *UserRequest) (*model.User, *model.Token, *app.ErrorJSON) {
	userModel := model.NewUser()
	fillModel(userModel, req)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, &app.ErrorJSON{
			Status:  http.StatusInternalServerError,
			Message: lang.M(ctx.Lang(), "user.service.generate-password"),
			Error:   err,
		}
	}
	userModel.Password = string(hashedPassword)

	if _, err := app.Mongo.Create(userModel); err != nil {
		return nil, nil, &app.ErrorJSON{
			Message: lang.M(ctx.Lang(), "app.service.store"),
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	tokenModel, er := tokenService(ctx, userModel)
	if er != nil {
		return nil, nil, er
	}

	return userModel, tokenModel, nil
}

func tokenService(ctx *app.ControllerContext, userModel *model.User) (*model.Token, *app.ErrorJSON) {
	tokenModel := model.NewToken()
	tokenModel.UserID = userModel.ID
	if _, err := app.Mongo.Create(tokenModel); err != nil {
		return nil, &app.ErrorJSON{
			Message: lang.M(ctx.Lang(), "app.service.store"),
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return tokenModel, nil
}

func updateService(ctx *app.ControllerContext, req *UserRequest) (*model.User, *app.ErrorJSON) {
	entity := &model.User{}
	if err := app.Mongo.FindByHexID(entity, req.ID); err != nil {
		return nil, &app.ErrorJSON{
			Status:  http.StatusNotFound,
			Message: lang.M(ctx.Lang(), "app.not-found"),
			Error:   err.Error(),
		}
	}

	if ctx.User.ID != entity.ID {
		return nil, &app.ErrorJSON{
			Status:  http.StatusUnauthorized,
			Message: lang.M(ctx.Lang(), "app.unautorized"),
			Error:   lang.M(ctx.Lang(), "app.unautorized"),
		}
	}

	fillModel(entity, req)
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, &app.ErrorJSON{
				Message: lang.M(ctx.Lang(), "user.service.generate-password"),
				Error:   err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		entity.Password = string(hashedPassword)
	}
	if _, err := app.Mongo.Update(entity); err != nil {
		return nil, &app.ErrorJSON{
			Message: lang.M(ctx.Lang(), "app.service.update"),
			Error:   err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return entity, nil
}

func loginService(ctx *app.ControllerContext, req *LoginRequest) (*model.User, *model.Token, *app.ErrorJSON) {
	entity := &model.User{}
	if validate.Email(req.User) {
		if err := app.Mongo.FindOneByField(entity, "email", req.User); err != nil {
			return nil, nil, &app.ErrorJSON{
				Message: lang.M(ctx.Lang(), "user.service.unautorized"),
				Error:   lang.M(ctx.Lang(), "user.service.unautorized"),
				Status:  http.StatusUnauthorized,
			}
		}
	} else {
		if err := app.Mongo.FindOneByField(entity, "phone", req.User); err != nil {
			return nil, nil, &app.ErrorJSON{
				Message: lang.M(ctx.Lang(), "user.service.unautorized"),
				Error:   lang.M(ctx.Lang(), "user.service.unautorized"),
				Status:  http.StatusUnauthorized,
			}
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(req.Password))
	if err != nil {
		return nil, nil, &app.ErrorJSON{
			Message: lang.M(ctx.Lang(), "user.service.unautorized"),
			Error:   lang.M(ctx.Lang(), "user.service.unautorized"),
			Status:  http.StatusUnauthorized,
		}
	}

	tokenModel, er := tokenService(ctx, entity)
	if er != nil {
		return nil, nil, er
	}

	return entity, tokenModel, nil
}

func deleteService(ctx *app.ControllerContext) *app.ErrorJSON {
	entity := &model.User{}
	id := ctx.Request.URL.Query().Get("u")
	if id == "" {
		return &app.ErrorJSON{
			Status:  http.StatusBadRequest,
			Message: lang.M(ctx.Lang(), "app.bad-request"),
			Error:   lang.M(ctx.Lang(), "app.request.query-params"),
		}
	}
	if err := app.Mongo.FindByHexID(entity, id); err != nil {
		return &app.ErrorJSON{
			Status:  http.StatusNotFound,
			Message: lang.M(ctx.Lang(), "app.not-found"),
			Error:   err.Error(),
		}
	}

	// if ctx.User.ID != entity.ID {
	// 	return &app.ErrorJSON{
	// 		Status:  http.StatusUnauthorized,
	// 		Message: lang.M(ctx.Lang(), "app.unautorized"),
	// 		Error:   lang.M(ctx.Lang(), "app.unautorized"),
	// 	}
	// }

	if _, err := app.Mongo.Delete(entity); err != nil {
		return &app.ErrorJSON{
			Status:  http.StatusInternalServerError,
			Message: lang.M(ctx.Lang(), "app.service.delete"),
			Error:   err.Error(),
		}
	}
	return nil
}

func fillModel(m *model.User, req *UserRequest) {

	m.Name = req.Name
	if req.Email != "" {
		m.Email = &req.Email
	}
	if req.Phone != "" {
		m.Phone = &req.Phone
	}
	m.UpdatedAt = time.Now()
}
