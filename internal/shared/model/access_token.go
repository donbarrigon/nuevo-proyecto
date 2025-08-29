package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	. "github.com/donbarrigon/nuevo-proyecto/internal/app/qb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AccessToken struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      bson.ObjectID `bson:"user_id"       json:"user_id"`
	User        *User         `bson:"user,omitempty" json:"user,omitempty"`
	Token       string        `bson:"token"         json:"token"`
	Permissions []string      `bson:"permissions"   json:"permissions"`
	CreatedAt   time.Time     `bson:"created_at"    json:"created_at"`
	ExpiresAt   time.Time     `bson:"expires_at"    json:"expires_at"`
	app.Odm     `bson:"-" json:"-"`
}

func (t *AccessToken) CollectionName() string { return "access_tokens" }
func (t *AccessToken) GetID() bson.ObjectID   { return t.ID }
func (t *AccessToken) SetID(id bson.ObjectID) { t.ID = id }

func (t *AccessToken) BeforeCreate() app.Error {
	t.CreatedAt = time.Now()
	t.ExpiresAt = t.generateExpiresAt()
	return nil
}

func (t *AccessToken) BeforeUpdate() app.Error {
	t.ExpiresAt = t.generateExpiresAt()
	return nil
}

func NewAccessToken() *AccessToken {
	token := &AccessToken{}
	token.Odm.Model = token
	return token
}

func (t *AccessToken) Generate(userID bson.ObjectID, permissions []string) app.Error {

	t.UserID = userID
	t.Token = t.generateToken()
	t.Permissions = permissions
	// t.CreatedAt = time.Now()
	// t.ExpiresAt = t.generateExpiresAt()

	return t.Create()
}

func (t *AccessToken) Refresh() app.Error {
	// t.Token = t.generateToken()
	// t.ExpiresAt = t.generateExpiresAt()
	return t.UpdateOne(Filter(Where("_id", Eq(t.ID))),
		Set(Element("expires_at", t.generateExpiresAt())),
	)
}

func (t *AccessToken) generateToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		app.PrintWarning("Fail to create access token: " + err.Error())
	}
	return hex.EncodeToString(bytes)
}

func (t *AccessToken) generateExpiresAt() time.Time {
	return time.Now().Add(time.Duration(app.Env.SESSION_DURATION) * time.Minute)
}

func (t *AccessToken) Can(permissionNames ...string) app.Error {
	// var result struct {
	// 	ExpiresAt time.Time `bson:"expires_at"`
	// }
	// err := app.DB.Collection(t.CollectionName()).FindOne(context.TODO(), bson.D{
	// 	{Key: "_id", Value: t.ID},
	// 	{Key: "permissions", Value: bson.D{{Key: "$in", Value: permissionNames}}},
	// },
	// 	options.FindOne().SetProjection(bson.M{"expires_at": 1}),
	// ).Decode(&result)
	// if err != nil {
	// 	return app.Errors.Forbidden(err)
	// }
	// if result.ExpiresAt.Before(time.Now()) {
	// 	return app.Errors.Forbiddenf("access denied: token expired at :expires_at", app.Entry{Key: "expires_at", Value: result.ExpiresAt})
	// }
	// return nil
	for _, permission := range t.Permissions {
		for _, permissionName := range permissionNames {
			if permission == permissionName {
				return nil
			}
		}
	}
	return app.Errors.Forbiddenf("access denied: missing permission: :permission", app.Entry{Key: "permission", Value: permissionNames})
}

func (t *AccessToken) GetUserID() bson.ObjectID {
	return t.UserID
}

func (t *AccessToken) HasRole(roleName ...string) app.Error {
	return t.User.HasRole(roleName...)
}

func (t *AccessToken) Anonymous() *AccessToken {
	var id bson.ObjectID // zero value: "000000000000000000000000"
	var timeZero time.Time
	return &AccessToken{
		ID:        id,
		UserID:    id,
		Token:     id.Hex(),
		CreatedAt: timeZero,
		ExpiresAt: timeZero,
	}
}
