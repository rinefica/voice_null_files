package cli

// RegisterRequestModel структура для регистрации пользователя.
type RegisterRequestModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
