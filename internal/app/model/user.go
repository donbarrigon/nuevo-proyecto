package model

import (
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
	"github.com/donbarrigon/nuevo-proyecto/pkg/validate"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"-"`
	Name      string        `bson:"name" json:"name"`
	Email     string        `bson:"email,omitempty" json:"email"`
	Phone     string        `bson:"phone,omitempty" json:"phone"`
	Password  string        `bson:"password" json:"password"`
	Tokens    *[]Token      `bson:"tokens" json:"tokens"`
	CreatedAt time.Time     `bson:"created_at" json:"-"`
	UpdatedAt time.Time     `bson:"updated_at" json:"-"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty" json:"-"`
}

func NewUser() *User {
	return &User{
		ID:        bson.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) CollectionName() string {
	return "users"
}

func (u *User) GetID() bson.ObjectID {
	return u.ID
}

func (u *User) SetID(id bson.ObjectID) {
	u.ID = id
}

func (u *User) Default() {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(l string) errors.Error {
	err := &errors.Err{}

	if strings.TrimSpace(u.Name) != "" {
		if len(u.Name) < 3 {
			err.Append("name", lang.TT(l, "Minimo %v caracteres", 3))
		}
		if len(u.Name) > 255 {
			err.Append("name", lang.TT(l, "Maximo %v caracteres", 255))
		}
	} else {
		err.Append("name", lang.TT(l, "Este campo es requerido"))
	}

	if strings.TrimSpace(u.Email) != "" {
		if len(u.Email) > 255 {
			err.Append("email", lang.TT(l, "Maximo %v caracteres", 255))
		}
		if !validate.Email(u.Email) {
			err.Append("email", lang.TT(l, "El formato de email es invalido"))
		}
	}

	if strings.TrimSpace(u.Phone) != "" {
		if len(u.Phone) < 5 {
			err.Append("phone", lang.TT(l, "Minimo %v caracteres", 5))
		}
		if len(u.Phone) > 255 {
			err.Append("phone", lang.TT(l, "Maximo %v caracteres", 255))
		}
	}

	if strings.TrimSpace(u.Password) != "" {
		if len(u.Password) < 8 {
			err.Append("password", lang.TT(l, "Minimo %v caracteres", 8))
		}
		if len(u.Password) > 32 {
			err.Append("password", lang.TT(l, "Maximo %v caracteres", 32))
		}
	} else {
		err.Append("password", lang.TT(l, "Este campo es requerido"))
	}

	if strings.TrimSpace(u.Email) == "" && strings.TrimSpace(u.Phone) == "" {
		err.Append("email", lang.TT(l, "Almenos uno es requerido"))
		err.Append("email", lang.TT(l, "Almenos uno es requerido"))
	}

	return err.Errors()
}
