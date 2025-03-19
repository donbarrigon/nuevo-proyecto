package user

import (
	"errors"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/model"
	"github.com/donbarrigon/nuevo-proyecto/pkg/validate"
	"golang.org/x/crypto/bcrypt"
)

func storeService(req *UserRequest) (*model.User, error) {
	userModel := model.NewUser()
	fillUserModel(userModel, req)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userModel.Password = string(hashedPassword)

	if _, err := app.DB.Create(userModel); err != nil {
		return nil, err
	}
	return userModel, nil
}

func updateService(req *UserRequest, authUser *model.User) (*model.User, error) {
	userModel := &model.User{}
	if err := app.DB.FindByHexID(userModel, req.ID); err != nil {
		return nil, err
	}

	if authUser.ID != userModel.ID {
		return nil, errors.New("No tienes permisos para modificar este usuario")
	}

	fillUserModel(userModel, req)
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		userModel.Password = string(hashedPassword)
	}
	if _, err := app.DB.Update(userModel); err != nil {
		return nil, err
	}
	return userModel, nil
}

func tokenStoreService(userModel *model.User) (*model.Token, error) {
	tokenModel := model.NewToken()
	tokenModel.UserID = userModel.ID
	if _, err := app.DB.Create(tokenModel); err != nil {
		return nil, err
	}
	return tokenModel, nil
}

func loginService(req *LoginRequest) *model.User {
	userModel := &model.User{}
	if validate.Email(req.User) {

	}
	if err := app.DB.FindOneByField(userModel, "email", req.User); err != nil {
		return nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(req.Password))
	if err != nil {
		return nil
	}

	return userModel
}

func fillUserModel(userModel *model.User, req *UserRequest) {

	userModel.Name = req.Name
	if req.Email != "" {
		userModel.Email = &req.Email
	}
	if req.Phone != "" {
		userModel.Phone = &req.Phone
	}
	userModel.UpdatedAt = time.Now()
}
