package validator

type StorePermission struct {
	Name string `json:"name" rules:"required|alpha_spaces|max:255"`
}
