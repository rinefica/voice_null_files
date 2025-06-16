package info_data

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rinefica/voice_null_files/internal/domain/model"
	"github.com/rinefica/voice_null_files/internal/services/crypto"
	"github.com/rinefica/voice_null_files/internal/storage"
)

// InfoDataService сервис для сохранения текстовой информации.
type InfoDataService interface {
	SaveInfoData(c *gin.Context)
	InfoData(c *gin.Context)
}

type InfoDataServiceImpl struct {
	log      *slog.Logger
	saver    storage.InfoDataSaver
	provider storage.InfoData
	crypto   crypto.CryptoService
}

func NewInfoDataService(log *slog.Logger, saver storage.InfoDataSaver, provider storage.InfoData, crypto crypto.CryptoService) InfoDataService {
	return &InfoDataServiceImpl{
		log:      log,
		saver:    saver,
		provider: provider,
		crypto:   crypto,
	}
}

// SaveInfoData сохраняет данные зарегистрированного пользователя.
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
	cryptoData, err := s.crypto.Encrypt([]byte(requestBody.Data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't save data"})
	}
	if err := s.saver.SaveInfoData(
		ctx,
		cryptoData,
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

// InfoData получает данные по ключу с проверкой принадлежности пользователю.
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
	decryptData, err := s.crypto.Decrypt(infoDataModel.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	infoDataModel.Data = string(decryptData)
	c.JSON(http.StatusOK, infoDataModel)
}
