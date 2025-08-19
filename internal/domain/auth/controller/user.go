package controller

import (
	"net/http"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	. "github.com/donbarrigon/nuevo-proyecto/internal/app/qb"
	"github.com/donbarrigon/nuevo-proyecto/internal/domain/auth/policy"
	"github.com/donbarrigon/nuevo-proyecto/internal/domain/auth/resource"
	"github.com/donbarrigon/nuevo-proyecto/internal/domain/auth/validator"
	"github.com/donbarrigon/nuevo-proyecto/internal/shared/model"
	"github.com/donbarrigon/nuevo-proyecto/internal/shared/service"
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

func UserTrashed(ctx *app.HttpContext) {
	if err := policy.UserDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	users := []*model.User{}
	if err := user.Find(users, Document(OnlyTrashed())); err != nil {
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

	req := &validator.StoreUser{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	user.Email = req.Email
	app.Fill(user.Profile, req)

	hashedPassword, er := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
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
		app.Log.Warning("User role does not exist. Run the seed command to populate initial data.", app.E("error", err))
	} else {
		user.RoleIDs = []bson.ObjectID{role.ID}
	}

	if err := user.Create(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.ActivityRecord(user.ID.Hex(), user, "create", user)
	go service.SendVerificationEmail(user)

	runLogin(ctx, req.Email, req.Password)
}

func Login(ctx *app.HttpContext) {

	req := &validator.UserLogin{}
	if err := ctx.GetBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	runLogin(ctx, req.Email, req.Password)
}
func runLogin(ctx *app.HttpContext, email string, password string) {

	user := model.NewUser()
	err := user.AggregateOne(mongo.Pipeline{
		Match(Where("email", Eq(email))),
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
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
	go model.ActivityRecord(user.ID.Hex(), accessToken, "login")

	ctx.ResponseOk(resource.NewUserLogin(user, accessToken))

}

func UserUpdateEmail(ctx *app.HttpContext) {

	req := &validator.UpdateUserEmail{}
	if err := ctx.ValidateBody(req); err != nil {
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.ResponseError(&app.Err{
			Status:  http.StatusUnauthorized,
			Message: "Invalid login credentials.",
			Err:     "Invalid login credentials.",
		})
		return
	}

	oldEmail := user.Email
	user.Email = req.Email
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

	req := &validator.UpdateUserProfile{}
	if err := ctx.ValidateBody(req); err != nil {
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

	dirty, err := app.FillDirty(user.Profile, req)
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

	req := &validator.UpdateUserPassword{}
	if err := ctx.ValidateBody(req); err != nil {
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.ResponseError(&app.Err{
			Status:  http.StatusUnauthorized,
			Message: "Invalid login credentials.",
			Err:     "Invalid login credentials.",
		})
		return
	}

	hashedPassword, er := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
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

func UserConfirmEmail(ctx *app.HttpContext) {

	verificationCode := model.NewVerificationCode()
	if err := verificationCode.FindOne(Document(Where("code", Eq(ctx.Params["code"])))); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Document(Where("_id", Eq(verificationCode.UserID)))); err != nil {
		ctx.ResponseError(err)
		return
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := verificationCode.Delete(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.ActivityRecord(user.ID.Hex(), user, "update", map[string]any{"email_verified_at": user.EmailVerifiedAt})

	ctx.ResponseNoContent()
}

func UserRevertEmail(ctx *app.HttpContext) {
	verificationCode := model.NewVerificationCode()
	if err := verificationCode.FindOne(Document(Where("code", Eq(ctx.Params["code"])))); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Document(Where("_id", Eq(verificationCode.UserID)))); err != nil {
		ctx.ResponseError(err)
		return
	}

	if verificationCode.Metadata["email"] == "" {
		ctx.ResponseError(&app.Err{
			Status:  http.StatusBadRequest,
			Message: "Invalid verification code: metadata email.",
			Err:     "Invalid verification code: metadata email.",
		})
		return
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	user.Email = verificationCode.Metadata["email"]
	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := verificationCode.Delete(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.ActivityRecord(user.ID.Hex(), user, "update", map[string]any{"email": user.Email, "email_verified_at": user.EmailVerifiedAt})

	ctx.ResponseNoContent()
}

func UserDestroy(ctx *app.HttpContext) {

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

	if err := policy.UserDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := user.SoftDelete(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.ActivityRecord(ctx.Auth.UserID(), user, "soft-delete", nil)

	ctx.ResponseNoContent()
}

func UserRestore(ctx *app.HttpContext) {
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

	if err := policy.UserDelete(ctx); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := user.Restore(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.ActivityRecord(ctx.Auth.UserID(), user, "restore", nil)

	ctx.ResponseNoContent()
}

func Logout(ctx *app.HttpContext) {

	accessToken, ok := ctx.Auth.Token.(*model.AccessToken)
	if !ok {
		ctx.ResponseError(app.Errors.InternalServerErrorF("Invalid token type:"))
		return
	}

	if err := accessToken.Delete(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go model.ActivityRecord(ctx.Auth.UserID(), accessToken, "logout", accessToken)

	ctx.ResponseNoContent()
}
