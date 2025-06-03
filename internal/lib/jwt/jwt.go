package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/rinefica/voice_null_files/internal/domain/model"
	"time"
)

func CreateToken(user *model.User, tokenTTL time.Duration, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(tokenTTL).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
