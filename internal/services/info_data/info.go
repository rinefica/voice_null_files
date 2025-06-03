package info_data

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rinefica/voice_null_files/internal/domain/model"
	"github.com/rinefica/voice_null_files/internal/storage"
	"log/slog"
	"net/http"
	"time"
)

type InfoDataService interface {
	SaveInfoData(c *gin.Context)
	InfoData(c *gin.Context)
}

type InfoDataServiceImpl struct {
	log      *slog.Logger
	saver    storage.InfoDataSaver
	provider storage.InfoData
}

func NewInfoDataService(log *slog.Logger, saver storage.InfoDataSaver, provider storage.InfoData) InfoDataService {
	return &InfoDataServiceImpl{
		log:      log,
		saver:    saver,
		provider: provider,
	}
}

func (s *InfoDataServiceImpl) SaveInfoData(c *gin.Context) {
	tag := "SaveInfoData"
	log := s.log.With("tag", tag)

	userID := c.GetInt("user_id")
	log.Info("User ID: ", userID)

	var requestBody model.InfoDataBody
	if err := c.ShouldBind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if requestBody.Data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data is empty"})
	}

	if requestBody.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type is empty"})
	}

	uuid := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.saver.SaveInfoData(
		ctx,
		requestBody.Data,
		requestBody.Type,
		requestBody.AdditionalData,
		uuid,
		int64(userID),
	); err != nil {
		log.Error("save info data", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("save info data", uuid)
	c.JSON(http.StatusOK, gin.H{"info": uuid})
}

func (s *InfoDataServiceImpl) InfoData(c *gin.Context) {
	tag := "InfoData"
	log := s.log.With("tag", tag)

	userID := c.GetInt("user_id")
	log.Info("User ID: ", userID)

	uuid := c.Param("uuid")
	log.Info("UUID: ", uuid)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	infoDataModel, err := s.provider.InfoData(ctx, uuid, int64(userID))
	if err != nil {
		log.Error("Info data not found" + err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": infoDataModel,
	})
}
