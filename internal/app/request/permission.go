package request

type Permission struct {
	Name string `json:"name" rules:"required|alpha_spaces|max:255"`
}
