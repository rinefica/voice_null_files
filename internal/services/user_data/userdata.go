package user_data

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rinefica/voice_null_files/internal/storage"
	"log/slog"
	"net/http"
	"time"
)

type UserDataService interface {
	UserData(c *gin.Context)
}

type UserDataServiceImpl struct {
	log      *slog.Logger
	provider storage.UserData
}

func NewUserDataServiceImpl(log *slog.Logger, provider storage.UserData) *UserDataServiceImpl {
	return &UserDataServiceImpl{
		log:      log,
		provider: provider,
	}
}

func (s *UserDataServiceImpl) UserData(c *gin.Context) {
	tag := "AllUserData"
	log := s.log.With("tag", tag)

	userID := c.GetInt("user_id")
	log.Info("User ID: ", userID)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	data, err := s.provider.AllData(ctx, int64(userID))
	if err != nil {
		log.Error("get all data ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
