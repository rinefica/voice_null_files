package cli

// LoginRequestModel структура для авторизации пользователя.
type LoginRequestModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponseModel структура для ответа авторизации, из которой берем токен.
type LoginResponseModel struct {
	Token string `json:"token"`
}
