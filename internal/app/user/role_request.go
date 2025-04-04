package user

import (
	"github.com/donbarrigon/nuevo-proyecto/pkg/lang"
)

type RoleRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (request *RoleRequest) Validate(language string) map[string][]string {
	errMap := make(map[string][]string, 0)
	errFields := make([]string, 0)

	if request.Name == "" {
		errFields = append(errFields, lang.M(language, "app.request.required"))
	}

	if len(request.Name) > 255 {
		errFields = append(errFields, lang.M(language, "app.request.max.txt", 255))
	}

	if len(errFields) > 0 {
		errMap["name"] = errFields
	}

	return errMap
}
