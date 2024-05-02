package dto

type CreateProductInput struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    bool   `json:"is_admin"`
}

type UpdateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    bool   `json:"is_admin"`
}

type JWTInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
