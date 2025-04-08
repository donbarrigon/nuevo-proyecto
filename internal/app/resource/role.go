package resource

type Role struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Permissions []*Permission `json:"permission"`
}
