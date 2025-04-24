package request

import (
	"github.com/donbarrigon/nuevo-proyecto/pkg/errors"
)

type Permission struct {
	Name string `json:"name"`
}

func (p *Permission) Validate(l string) errors.Error {
	err := &errors.Err{}

	err.Append("name", MaxString(l, p.Name, 255))
	err.Append("name", Required(l, p.Name))

	p.Name.IsZero()

	return err.Errors()
}
