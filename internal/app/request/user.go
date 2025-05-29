package request

type StoreUser struct {
	Name                 string `json:"name" rules:"required|max:255"`
	Email                string `json:"email" rules:"required_without:phone|email|max:255"`
	Phone                string `json:"phone" rules:"required_without:email|max:255"`
	Password             string `json:"password" rules:"reqired|confirmed|min:8|max:32"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type UpdateUser struct {
	Name  string `json:"name" rules:"required|max:255"`
	Email string `json:"email" rules:"required_without:phone|email|max:255"`
	Phone string `json:"phone" rules:"required_without:email|max:255"`
}
