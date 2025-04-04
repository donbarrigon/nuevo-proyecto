package model

import (
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/core"
	"github.com/donbarrigon/nuevo-proyecto/pkg/validate"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id" json:"-"`
	Name      string        `bson:"name" json:"name"`
	Email     string        `bson:"email,omitempty" json:"email"`
	Phone     string        `bson:"phone,omitempty" json:"phone"`
	Password  string        `bson:"password" json:"password"`
	CreatedAt time.Time     `bson:"createdAt" json:"-"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"-"`
	DeletedAt time.Time     `bson:"deletedAt,omitempty" json:"-"`
}

func NewUser() *User {
	return &User{
		ID:        bson.NewObjectID(),
		CreatedAt: time.Now(),
	}
}

func (u *User) CollectionName() string {
	return "users"
}

func (u *User) GetID() bson.ObjectID {
	return u.ID
}

func (u *User) Validate(ctx *core.Context) core.Error {
	errMap := make(map[string][]string, 0)
	errFields := make([]string, 0)

	if strings.TrimSpace(u.Name) != "" {
		if len(u.Name) < 3 {
			errFields = append(errFields, ctx.TT("Minimo %v caracteres", 3))
		}
		if len(u.Name) > 255 {
			errFields = append(errFields, ctx.TT("Maximo %v caracteres", 255))
		}
		if len(errFields) > 0 {
			errMap["name"] = errFields
			errFields = make([]string, 0)
		}
	} else {

	}

	if strings.TrimSpace(u.Email) != "" {
		if len(u.Email) > 255 {
			errFields = append(errFields, ctx.TT("Maximo %v caracteres", 255))
		}
		if !validate.Email(u.Email) {
			errFields = append(errFields, ctx.TT("El formato de email es invalido"))
		}
		if len(errFields) > 0 {
			errMap["email"] = errFields
			errFields = make([]string, 0)
		}
	}

	if strings.TrimSpace(u.Phone) != "" {
		if len(u.Phone) < 5 {
			errFields = append(errFields, ctx.TT("Minimo %v caracteres", 5))
		}
		if len(u.Phone) > 255 {
			errFields = append(errFields, ctx.TT("Maximo %v caracteres", 255))
		}
		if len(errFields) > 0 {
			errMap["phone"] = errFields
			errFields = make([]string, 0)
		}
	}

	if strings.TrimSpace(u.Password) != "" {
		if len(u.Password) < 8 {
			errFields = append(errFields, ctx.TT("Minimo %v caracteres", 8))
		}
		if len(u.Password) > 32 {
			errFields = append(errFields, ctx.TT("Maximo %v caracteres", 32))
		}
		if len(errFields) > 0 {
			errMap["password"] = errFields
		}
	} else {
		errMap["password"] = []string{ctx.TT("Este campo es requerido")}
	}

	if strings.TrimSpace(u.Email) == "" && strings.TrimSpace(u.Phone) == "" {
		errMap["email"] = []string{ctx.TT("Este campo es requerido")}
		errMap["phone"] = []string{ctx.TT("Este campo es requerido")}
	}
	return nil
}
