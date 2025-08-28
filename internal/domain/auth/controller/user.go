package controller

import (
	"crypto/rand"
	"math/big"
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

	if user.Profile.CityID, er = bson.ObjectIDFromHex(req.CityID); er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	role := model.NewRole()
	if err := role.FindOne(Filter(Where("name", Eq("user")))); err != nil {
		app.PrintWarning("User role does not exist. Run the seed command to populate initial data.", app.E("error", err))
	} else {
		user.RoleIDs = []bson.ObjectID{role.ID}
	}

	if err := user.Create(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.ActivityRecord(user.ID, user, "create", user)
	go service.SendEmailConfirm(user)

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
	accessToken := model.NewAccessToken()
	if err := accessToken.Generate(user.ID, permissions); err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.ActivityRecord(user.ID, accessToken, "login")

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
	if err := user.FindOne(Filter(Where("_id", Eq(id)))); err != nil {
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
	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := user.UpdateOne(
		Filter(Where("_id", Eq(user.ID))),
		Unset("email_verified_at"),
	); err != nil {
		if err.GetMessage() != "No records updated" {
			ctx.ResponseError(err)
			return
		}
	}
	// filter := bson.D{bson.E{Key: "_id", Value: o.Model.GetID()}}
	// update := bson.D{bson.E{Key: "$unset", Value: bson.D{{Key: "deleted_at", Value: nil}}}}

	go service.ActivityRecord(ctx.Auth.UserID(), user, "update-email", map[string]string{"email": user.Email})
	go service.SendEmailConfirm(user)
	go service.SendEmailChanged(user, oldEmail)

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
	if err := user.FindOne(Filter(Where("_id", Eq(id)))); err != nil {
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

	if user.Profile.CityID.Hex() != req.CityID {
		dirty["city_id"] = req.CityID
		if user.Profile.CityID, er = bson.ObjectIDFromHex(req.CityID); er != nil {
			ctx.ResponseError(app.Errors.HexID(er))
			return
		}
	}

	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.ActivityRecord(ctx.Auth.UserID(), user, "update-profile", dirty)

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
	if err := user.FindOne(Filter(Where("_id", Eq(id)))); err != nil {
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

	hashedPassword, er := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
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

	accesToken := model.NewAccessToken()
	if err := accesToken.DeleteMany(Filter(Where("user_id", Eq(user.ID)))); err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.ActivityRecord(ctx.Auth.UserID(), user, "update-password", map[string]string{"password": user.Password})
	go service.SendEmailPasswordChanged(user)

	ctx.ResponseOk(user)
}

func UserConfirmEmail(ctx *app.HttpContext) {

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	verificationCode := model.NewVerificationCode()
	if err := verificationCode.FindOne(Filter(
		Where("user_id", Eq(id)),
		Where("type", Eq("email-verification")),
		Where("code", Eq(ctx.Params["code"])),
	)); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Filter(Where("_id", Eq(verificationCode.UserID)))); err != nil {
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

	go service.ActivityRecord(user.ID, user, "confirm-email", map[string]any{"email_verified_at": user.EmailVerifiedAt})

	ctx.ResponseOk(map[string]string{"message": "Email verified.", "email_verified_at": user.EmailVerifiedAt.Format(time.RFC3339)})
}

func UserRevertEmail(ctx *app.HttpContext) {

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	verificationCode := model.NewVerificationCode()
	if err := verificationCode.FindOne(Filter(
		Where("user_id", Eq(id)),
		Where("type", Eq("email-change-revert")),
		Where("code", Eq(ctx.Params["code"])),
	)); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Filter(Where("_id", Eq(verificationCode.UserID)))); err != nil {
		ctx.ResponseError(err)
		return
	}

	if verificationCode.Metadata["old_email"] == "" {
		ctx.ResponseError(&app.Err{
			Status:  http.StatusBadRequest,
			Message: "Invalid verification code: metadata old_email.",
			Err:     "Invalid verification code: metadata old_email.",
		})
		return
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	user.Email = verificationCode.Metadata["old_email"]
	if err := user.Update(); err != nil {
		ctx.ResponseError(err)
		return
	}

	if err := verificationCode.Delete(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.ActivityRecord(user.ID, user, "revert-email", map[string]any{"email": user.Email, "email_verified_at": user.EmailVerifiedAt})

	ctx.ResponseOk(map[string]string{"message": "Email reverted.", "email": user.Email, "email_verified_at": user.EmailVerifiedAt.Format(time.RFC3339)})
}

func UserForgotPassword(ctx *app.HttpContext) {
	req := &validator.ForgotPassword{}
	if err := ctx.ValidateBody(req); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Filter(Where("email", Eq(req.Email)))); err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.SendEmailForgotPassword(user)

	ctx.ResponseOk(map[string]string{"message": "Email sent."})
}

func UserResetPassword(ctx *app.HttpContext) {
	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	verificationCode := model.NewVerificationCode()
	if err := verificationCode.FindOne(Filter(
		Where("user_id", Eq(id)),
		Where("type", Eq("reset-password")),
		Where("code", Eq(ctx.Params["code"])),
	)); err != nil {
		ctx.ResponseError(err)
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Filter(Where("_id", Eq(verificationCode.UserID)))); err != nil {
		ctx.ResponseError(err)
		return
	}

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 12)
	for i := range result {
		num, er := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if er != nil {
			ctx.ResponseError(app.Errors.InternalServerErrorf("Error generating password: :error", app.E("error", er.Error())))
		}
		result[i] = letters[num.Int64()]
	}
	newPassword := string(result)

	hashedPassword, er := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
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

	if err := verificationCode.Delete(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.ActivityRecord(user.ID, user, "reset-password", map[string]any{"password": newPassword})
	go service.SendMailNewPassword(user, newPassword)

	ctx.ResponseOk(map[string]string{"message": "Password reset."})
}

func UserDestroy(ctx *app.HttpContext) {

	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Filter(Where("_id", Eq(id)))); err != nil {
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

	go service.ActivityRecord(ctx.Auth.UserID(), user, "soft-delete", nil)

	ctx.ResponseNoContent()
}

func UserRestore(ctx *app.HttpContext) {
	id, er := bson.ObjectIDFromHex(ctx.Params["id"])
	if er != nil {
		ctx.ResponseError(app.Errors.HexID(er))
		return
	}

	user := model.NewUser()
	if err := user.FindOne(Filter(Where("_id", Eq(id)))); err != nil {
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

	go service.ActivityRecord(ctx.Auth.UserID(), user, "restore", nil)

	ctx.ResponseNoContent()
}

func Logout(ctx *app.HttpContext) {

	accessToken, ok := ctx.Auth.Token.(*model.AccessToken)
	if !ok {
		ctx.ResponseError(app.Errors.InternalServerErrorf("Invalid token type:"))
		return
	}

	if err := accessToken.Delete(); err != nil {
		ctx.ResponseError(err)
		return
	}

	go service.ActivityRecord(ctx.Auth.UserID(), accessToken, "logout", accessToken)

	ctx.ResponseNoContent()
}
