package main

import (
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rinefica/voice_null_files/internal/app"
	"github.com/rinefica/voice_null_files/internal/config"
	"github.com/rinefica/voice_null_files/internal/lib/sl"
)

// Запуск приложения сервера, чтение конфига, подключение к БД. Поддержка graceful shutdown.
func main() {
	cfg := config.MustLoad()

	log := sl.SetupLogger(cfg.Env)
	log.Debug("config env:: ", slog.Any("config", cfg))

	log.Info("version: ", slog.Any("version", cfg.Version))

	applctn := app.NewApp(
		log,
		cfg.Server.Port,
		cfg.StoragePath,
		cfg.TokenTTL,
		cfg.Secret,
		cfg.Key,
		cfg.Server.UseSSL,
	)

	go applctn.Server.MustRun(cfg.Server.UseSSL)

	db, err := sql.Open("postgres", cfg.StoragePath)
	if err != nil {
		log.Debug("DB connect Failed")
		log.Error(err.Error())
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Debug("DB Ping Failed")
		log.Error(err.Error())
	}
	log.Debug("DB Connection started successfully")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	applctn.Server.Stop()

	log.Debug("app stopped")
}
