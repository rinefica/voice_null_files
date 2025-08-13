package files

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rinefica/voice_null_files/internal/storage"
)

// FileService сервис для сохранения и получения файлов с сервера
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

// File метод для получения ранее загруженного файла на сервер,
// использует uuid файла для поиска и user_id для проверки доступности пользователю.
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

// SaveFile сохраняет файл на сервер, генерирует uuid и сохраняет с этим именем,
// также сохраняет в БД принадлежность файла пользователю по user_id.
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

	c.JSON(http.StatusOK, fmt.Sprintf("'%s' uploaded! uuid %s", f.Filename, filename))
}
