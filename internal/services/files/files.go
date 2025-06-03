package files

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rinefica/voice_null_files/internal/storage"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FileService interface {
	SaveFile(c *gin.Context)
	File(c *gin.Context)
}

type FileServiceImpl struct {
	log          *slog.Logger
	fileSaver    storage.FileSaver
	fileProvider storage.File
}

func NewFileService(log *slog.Logger, saver storage.FileSaver, fileProvider storage.File) FileService {
	return &FileServiceImpl{
		log:          log,
		fileSaver:    saver,
		fileProvider: fileProvider,
	}
}

func (s *FileServiceImpl) File(c *gin.Context) {
	tag := "File"
	log := s.log.With("tag", tag)

	userID := c.GetInt("user_id")
	log.Info("User ID: ", userID)

	uuid := c.Param("uuid")
	log.Info("UUID: ", uuid)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	file, err := s.fileProvider.File(ctx, uuid, int64(userID))
	if err != nil {
		log.Info("GetUploadedFile")
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	path := filepath.Join(".", "public/", file.UUID)
	log.Info("file path: ", path)
	c.FileAttachment(path, file.Filename)
}

func (s *FileServiceImpl) SaveFile(c *gin.Context) {
	tag := "SaveFile"
	log := s.log.With("tag", tag)

	userID := c.GetInt("user_id")
	log.Info("User ID: ", userID)

	form, err := c.MultipartForm()
	log.Info("form == ", form)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	files := form.File["file"]

	f := files[0]
	log.Info("get file", f.Filename, f.Size)
	filename := uuid.New().String()
	log.Info("save file", filename)

	path := filepath.Join(".", "public/")
	log.Info("path = ", path)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Info("MkdirAll")
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := c.SaveUploadedFile(f, filepath.Join(path, filename)); err != nil {
		log.Info("SaveUploadedFile")
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	log.Info("success save file")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.fileSaver.SaveFile(ctx, f.Filename, filename, int64(userID)); err != nil {
		log.Info("fileSaver.SaveFile")
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	log.Info("success save file to db")

	c.JSON(http.StatusOK, fmt.Sprintf("'%s' uploaded!", f.Filename))
}
