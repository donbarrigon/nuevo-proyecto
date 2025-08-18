package controller

import (
	"net/http"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	. "github.com/donbarrigon/nuevo-proyecto/internal/app/qb"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/policy"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/resource"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/validator"
	"github.com/donbarrigon/nuevo-proyecto/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func UserIndex(ctx *app.HttpContext) {
	if err := policy.UserViewAny(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	users := []*model.User{}
	if err := user.Find(users, Document(WithOutTrashed())); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseOk(users)
}

func UserShow(ctx *app.HttpContext) {
	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	user := model.NewUser()
	err := user.AggregateOne(mongo.Pipeline{
		Match(Where("_id", Eq(id))),
		user.WithRoles(),
		user.WhithPermissions(),
		user.WithTokens(),
	})
	if err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := policy.UserView(ctx, user); err != nil {
		ctx.ResponseError(err)
		return
	}

	ctx.ResponseOk(user)
}

func UserStore(ctx *app.HttpContext) {

	if err := policy.UserCreate(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	validator := &validator.StoreUser{}
	if err := ctx.ValidateBody(validator); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	user.Email = validator.Email
	app.Fill(user.Profile, validator)

	hashedPassword, er := bcrypt.GenerateFromPassword([]byte(validator.Password), bcrypt.DefaultCost)
	if er != nil {
		ctx.ResponseError(&app.Err{
			Status:  http.StatusInternalServerError,
			Message: "Password encryption failed",
			Err:     er.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	role := model.NewRole()
	if err := role.FindOne(Document(Where("name", Eq("user")))); err != nil {
		ctx.ResponseError(err)
		return
	}
	user.RoleIDs = []bson.ObjectID{role.ID}

	if err := user.Create(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.ActivityRecord(user.ID.Hex(), user, "create", user)
	go service.SendVerificationEmail(user)

	Login(ctx)
}

func Login(ctx *app.HttpContext) {

	validator := &validator.UserLogin{}
	if err := ctx.GetBody(validator); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	err := user.AggregateOne(mongo.Pipeline{
		Match(Where("email", Eq(validator.Email))),
		user.WithRoles(),
		user.WhithPermissions(),
	})
	if err != nil {
		ctx.ResponseError(&app.Err{
			Status:  http.StatusUnauthorized,
			Message: "Invalid login credentials.",
			Err:     "Invalid login credentials.",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(validator.Password)); err != nil {
		ctx.ResponseError(&app.Err{
			Status:  http.StatusUnauthorized,
			Message: "Invalid login credentials.",
			Err:     "Invalid login credentials.",
		})
		return
	}

	permissions := []string{}
	exists := false
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {

			for _, p := range permissions {
				if p == permission.Name {
					exists = true
					break
				}
			}
			if !exists {
				permissions = append(permissions, permission.Name)
			}
		}
	}

	for _, permission := range user.Permissions {
		for _, p := range permissions {
			if p == permission.Name {
				exists = true
				break
			}
		}
		if !exists {
			permissions = append(permissions, permission.Name)
		}
	}

	accessToken, err := model.NewAccessToken(user.ID, permissions)
	if err != nil {
		ctx.ResponseError(err)
		return
	}
	go model.ActivityRecord(user.ID.Hex(), accessToken, "login", accessToken)

	ctx.ResponseOk(resource.NewUserLogin(user, accessToken))

}

func UserUpdateEmail(ctx *app.HttpContext) {

	validator := &validator.UpdateUserEmail{}
	if err := ctx.ValidateBody(validator); err != nil {
		ctx.ResponseError(err)
		return
	}

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Document(Where("_id", Eq(id)))); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := policy.UserUpdate(ctx, user); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(validator.Password)); err != nil {
		ctx.ResponseError(&app.Err{
			Status:  http.StatusUnauthorized,
			Message: "Invalid login credentials.",
			Err:     "Invalid login credentials.",
		})
		return
	}

	oldEmail := user.Email
	user.Email = validator.Email
	user.EmailVerifiedAt = nil
	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.ActivityRecord(ctx.Auth.UserID(), user, "update", map[string]string{"email": user.Email})
	go service.SendVerificationEmail(user)
	go service.SendEmailChangeNotification(user, oldEmail)

	ctx.ResponseOk(user)
}

func UserUpdateProfile(ctx *app.HttpContext) {

	validator := &validator.UpdateUserProfile{}
	if err := ctx.ValidateBody(validator); err != nil {
		ctx.ResponseError(err)
		return
	}

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Document(Where("_id", Eq(id)))); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := policy.UserUpdate(ctx, user); err != nil {
		ctx.ResponseError(err)
		return
	}

	dirty, err := app.FillDirty(user.Profile, validator)
	if err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.ActivityRecord(ctx.Auth.UserID(), user, "update", dirty)

	ctx.ResponseOk(user)
}

func UserUpdatePassword(ctx *app.HttpContext) {

	validator := &validator.UpdateUserPassword{}
	if err := ctx.ValidateBody(validator); err != nil {
		ctx.ResponseError(err)
		return
	}

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Document(Where("_id", Eq(id)))); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := policy.UserUpdate(ctx, user); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(validator.Password)); err != nil {
		ctx.ResponseError(&app.Err{
			Status:  http.StatusUnauthorized,
			Message: "Invalid login credentials.",
			Err:     "Invalid login credentials.",
		})
		return
	}

	hashedPassword, er := bcrypt.GenerateFromPassword([]byte(validator.Password), bcrypt.DefaultCost)
	if er != nil {
		ctx.ResponseError(&app.Err{
			Status:  http.StatusInternalServerError,
			Message: "Password encryption failed",
			Err:     er.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.ActivityRecord(ctx.Auth.UserID(), user, "update", map[string]string{"password": user.Password})

	ctx.ResponseOk(user)
}
