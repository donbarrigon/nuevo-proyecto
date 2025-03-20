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

func storeService(ctx *app.HandlerContext, req *UserRequest) (*model.User, *app.ErrorJSON) {
	entity := model.NewUser()
	fillModel(entity, req)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &app.ErrorJSON{
			Status:  http.StatusInternalServerError,
			Message: lang.M(ctx.Lang(), "user.service.generate_password"),
			Error:   err,
		}
	}
	entity.Password = string(hashedPassword)

	if _, err := app.Mongo.Create(entity); err != nil {
		return nil, &app.ErrorJSON{
			Status:  http.StatusInternalServerError,
			Message: lang.M(ctx.Lang(), "app.service.store"),
			Error:   err,
		}
	}
	return entity, nil
}

func updateService(ctx *app.HandlerContext, req *UserRequest) (*model.User, *app.ErrorJSON) {
	entity := &model.User{}
	if err := app.Mongo.FindByHexID(entity, req.ID); err != nil {
		return nil, &app.ErrorJSON{
			Status:  http.StatusNotFound,
			Message: lang.M(ctx.Lang(), "app.not_found"),
			Error:   err.Error(),
		}
	}

	if ctx.User.ID != entity.ID {
		return nil, &app.ErrorJSON{
			Status:  http.StatusUnauthorized,
			Message: lang.M(ctx.Lang(), "app.unautorized"),
			Error:   lang.M(ctx.Lang(), "user.service.update.id"),
		}
	}

	fillModel(entity, req)
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, &app.ErrorJSON{
				Status:  http.StatusInternalServerError,
				Message: lang.M(ctx.Lang(), "user.service.generate_password"),
				Error:   err.Error(),
			}
		}
		entity.Password = string(hashedPassword)
	}
	if _, err := app.Mongo.Update(entity); err != nil {
		return nil, &app.ErrorJSON{
			Status:  http.StatusInternalServerError,
			Message: lang.M(ctx.Lang(), "app.service.update"),
			Error:   err.Error(),
		}
	}
	return entity, nil
}

func tokenStoreService(ctx *app.HandlerContext, userModel *model.User) (*model.Token, *app.ErrorJSON) {
	entity := model.NewToken()
	entity.UserID = userModel.ID
	if _, err := app.Mongo.Create(entity); err != nil {
		return nil, &app.ErrorJSON{
			Status:  http.StatusInternalServerError,
			Message: lang.M(ctx.Lang(), "app.service.store"),
			Error:   err.Error(),
		}
	}
	return entity, nil
}

func loginService(ctx *app.HandlerContext, req *LoginRequest) (*model.User, *app.ErrorJSON) {
	entity := &model.User{}
	if validate.Email(req.User) {
		if err := app.Mongo.FindOneByField(entity, "email", req.User); err != nil {
			return nil, &app.ErrorJSON{
				Status:  http.StatusUnauthorized,
				Message: lang.M(ctx.Lang(), "user.service.unautorized"),
				Error:   lang.M(ctx.Lang(), "user.service.unautorized"),
			}
		}
	} else {
		if err := app.Mongo.FindOneByField(entity, "phone", req.User); err != nil {
			return nil, &app.ErrorJSON{
				Status:  http.StatusUnauthorized,
				Message: lang.M(ctx.Lang(), "user.service.unautorized"),
				Error:   lang.M(ctx.Lang(), "user.service.unautorized"),
			}
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(entity.Password), []byte(req.Password))
	if err != nil {
		return nil, &app.ErrorJSON{
			Status:  http.StatusUnauthorized,
			Message: lang.M(ctx.Lang(), "user.service.unautorized"),
			Error:   lang.M(ctx.Lang(), "user.service.unautorized"),
		}
	}

	return entity, nil
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
