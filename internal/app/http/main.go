package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rinefica/voice_null_files/internal/services/auth"
	"github.com/rinefica/voice_null_files/internal/services/crypto"
	"github.com/rinefica/voice_null_files/internal/services/files"
	"github.com/rinefica/voice_null_files/internal/services/info_data"
	"github.com/rinefica/voice_null_files/internal/services/user_data"
	"github.com/rinefica/voice_null_files/internal/storage"
)

// App серверное приложение, использующее http протокол.
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
	key []byte,
	useSSL bool,
) *App {

	router := gin.Default()
	authService := auth.New(log, storage, storage, tokenTTL, secret)
	fileService := files.NewFileService(log, storage, storage)
	cryptoService := crypto.NewCryptoService(log, key)
	infoService := info_data.NewInfoDataService(log, storage, storage, cryptoService)
	userService := user_data.NewUserDataServiceImpl(log, storage)
	setupRouter(log, router, authService, fileService, infoService, userService, secret)

	config := &tls.Config{}
	if useSSL {
		cert, err := tls.LoadX509KeyPair("keys/server.crt", "keys/server.key")
		if err != nil {
			panic(err)
		}
		config = &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
		}
	}

	srv := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   router.Handler(),
		TLSConfig: config,
	}

	return &App{
		log:     log,
		router:  router,
		server:  srv,
		storage: storage,
	}
}

func (a *App) MustRun(useSSL bool) {
	if useSSL {
		if err := a.server.ListenAndServeTLS("keys/server.crt", "keys/server.key"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	} else {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}

func (a *App) Stop() {
	const tag = "httpApp.stop"
	log := a.log.With("tag", tag)

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

	// Авторизация в системе.
	router.POST("/login", authService.Login)
	// Регистрация в системе.
	router.POST("/register", authService.Register)

	// Проверка jwt-токена для всех последующих запросов.
	router.Use(authMiddleware(log, secret))

	// Получение ранее загруженного этим пользователем файла по uuid.
	router.GET("/api/file/:uuid", fileService.File)
	// Загрузка пользователем файла в систему.
	router.POST("/api/file/", fileService.SaveFile)

	// Получение текстовых данных пользователя по uuid.
	router.GET("/api/info_data/:uuid", infoService.InfoData)
	// Добавление текстовых данных в систему.
	router.POST("/api/info_data/", infoService.SaveInfoData)

	// Получение всех загруженных пользователем данных.
	router.GET("/api/info_data/all", userService.UserData)
}

// authMiddleware добавляет проверку на наличие jwt токена в запросе,
// для внутренних дальнейших обработок добавляет данные пользователя к контексту запроса.
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
