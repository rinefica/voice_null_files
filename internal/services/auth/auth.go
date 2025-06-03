package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rinefica/voice_null_files/internal/domain/model"
	"github.com/rinefica/voice_null_files/internal/lib/jwt"
	"github.com/rinefica/voice_null_files/internal/lib/sl"
	"github.com/rinefica/voice_null_files/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"time"
)

type AuthService interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}
type AuthServiceImpl struct {
	log          *slog.Logger
	userSaver    storage.UserSaver
	userProvider storage.UserProvider
	tokenTTL     time.Duration
	secret       string
}

func New(
	log *slog.Logger,
	userSaver storage.UserSaver,
	userProvider storage.UserProvider,
	tokenTTL time.Duration,
	secret string,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		tokenTTL:     tokenTTL,
		secret:       secret,
	}
}

func (a *AuthServiceImpl) Login(
	c *gin.Context,
) {
	var requestBody model.RegisterBody
	if err := c.ShouldBind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	email := requestBody.Email
	password := requestBody.Password
	const tag = "auth.login"
	log := a.log.With(slog.String("tag", tag))

	log.Info("start login")
	log.Info("body", requestBody)

	if email == "" || password == "" {
		log.Debug("credentials invalid email %s pswd %s", email, password)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	user, err := a.userProvider.User(c, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Debug("user not found", sl.Err(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		log.Error("failed get user", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Debug("invalid credentials", sl.Err(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	log.Info("success login")

	token, err := jwt.CreateToken(user, a.tokenTTL, a.secret)
	if err != nil {
		log.Error("failed create token", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *AuthServiceImpl) Register(
	c *gin.Context,
) {
	var requestBody model.RegisterBody
	if err := c.ShouldBind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	email := requestBody.Email
	password := requestBody.Password
	const tag = "auth.register"
	log := a.log.With(slog.String("tag", tag))

	log.Info("Register user %s", email)
	log.Info("body", requestBody)

	if email == "" || password == "" {
		log.Debug("credentials invalid email %s pswd %s", email, password)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate password hash"})
		return
	}

	uid, err := a.userSaver.SaveUser(c, email, passHash)
	if err != nil {
		log.Error("failed to save user", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}
	log.Debug("saved user %s", email)
	c.JSON(http.StatusOK, gin.H{"userId": uid})
}
