package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rinefica/voice_null_files/internal/services/auth"
	"github.com/rinefica/voice_null_files/internal/services/files"
	"github.com/rinefica/voice_null_files/internal/services/info_data"
	"github.com/rinefica/voice_null_files/internal/services/user_data"
	"github.com/rinefica/voice_null_files/internal/storage"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	log     *slog.Logger
	router  *gin.Engine
	server  *http.Server
	storage *storage.Storage
}

func NewApp(
	log *slog.Logger,
	port int,
	storage *storage.Storage,
	tokenTTL time.Duration,
	secret string,
) *App {

	router := gin.Default()
	authService := auth.New(log, storage, storage, tokenTTL, secret)
	fileService := files.NewFileService(log, storage, storage)
	infoService := info_data.NewInfoDataService(log, storage, storage)
	userService := user_data.NewUserDataServiceImpl(log, storage)
	setupRouter(log, router, authService, fileService, infoService, userService, secret)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.Handler(),
	}

	return &App{
		log:     log,
		router:  router,
		server:  srv,
		storage: storage,
	}
}

func (a *App) MustRun() {
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (a *App) Stop() {
	const tag = "httpApp.stop"
	log := a.log.With("tag", tag)
	log.Info("stopping grpc server")

	a.storage.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Error("Server Shutdown:", err)
	}
}

func setupRouter(
	log *slog.Logger,
	router *gin.Engine,
	authService auth.AuthService,
	fileService files.FileService,
	infoService info_data.InfoDataService,
	userService user_data.UserDataService,
	secret string) {

	// Route for generating tokens
	router.POST("/login", authService.Login)
	// Route for generating tokens
	router.POST("/register", authService.Register)

	// Middleware to check JWT on every request
	router.Use(authMiddleware(log, secret))
	router.GET("/api/file/:uuid", fileService.File)
	router.POST("/api/file/", fileService.SaveFile)

	router.GET("/api/info_data/:uuid", infoService.InfoData)
	router.POST("/api/info_data/", infoService.SaveInfoData)

	router.GET("/api/info_data/all", userService.UserData)
}

func authMiddleware(log *slog.Logger, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Stop further processing if unauthorized
			return
		}

		// Set the token claims to the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Info("claims: ", claims)
			c.Set("claims", claims)
			c.Set("user_id", int(claims["user_id"].(float64)))
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next() // Proceed to the next handler if authorized
	}
}
