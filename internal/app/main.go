package app

import (
	"github.com/rinefica/voice_null_files/internal/app/http"
	"github.com/rinefica/voice_null_files/internal/storage"

	"log/slog"
	"time"
)

type App struct {
	Server http.App
}

func NewApp(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	secret string,
) *App {

	strg, err := storage.NewStorage(log, storagePath)
	if err != nil {
		panic(err)
	}

	newApp := http.NewApp(log, grpcPort, strg, tokenTTL, secret)
	return &App{
		Server: *newApp,
	}
}
