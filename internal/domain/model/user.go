package model

// User структура для описания пользователя.
type User struct {
	// ID идентификатор в БД сервера.
	ID int64
	// Email почта, с которой была произведена регистрация.
	Email string
	// PasswordHash хеш пароля для проверок идентификации.
	PasswordHash []byte
}

// RegisterBody модель для регистрации и авторизации пользователя.
type RegisterBody struct {
	// Email почта, которую использует пользователь.
	Email string
	// Password пароль пользователя.
	Password string
}
