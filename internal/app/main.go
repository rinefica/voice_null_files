package app

import (
	"log/slog"
	"time"

	"github.com/rinefica/voice_null_files/internal/app/http"
	"github.com/rinefica/voice_null_files/internal/storage"
)

// App основное серверное приложение, в дальнейшем может расширяться различными реализациями и протоколами (например, добавить взаимодействие с клиентом по сокету, масштабирование системы на gRPC модули)
type App struct {
	Server http.App
}

func NewApp(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	secret string,
	key string,
	useSSL bool,
) *App {

	strg, err := storage.NewStorage(log, storagePath)
	if err != nil {
		panic(err)
	}

	newApp := http.NewApp(log, grpcPort, strg, tokenTTL, secret, []byte(key), useSSL)
	return &App{
		Server: *newApp,
	}
}
