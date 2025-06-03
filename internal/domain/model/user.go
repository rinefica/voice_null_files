package model

type User struct {
	ID           int64
	Email        string
	PasswordHash []byte
}

type RegisterBody struct {
	Email    string
	Password string
}
