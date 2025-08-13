package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config описывает конфигурацию сервера.
type Config struct {
	// Env описание текущей среды для возможности смены контекста (local, test, prod).
	Env string `yaml:"env"`
	// Version текущая версия приложения.
	Version string `yaml:"version"`
	// Secret ключ для работы с jwt токеном.
	Secret string `yaml:"secret"`
	// Key ключ для шифрования данных в БД.
	Key string `yaml:"key"`
	// StoragePath путь к БД.
	StoragePath string `yaml:"storage_path" env-required:"true"`
	// Server характеристики http(s)-сервера.
	Server ServerConfig `yaml:"server"`
	// TokenTTL время жизни токена доступа.
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
}

// ServerConfig структура для описания характеристик сервера.
type ServerConfig struct {
	// Port порт для http(s)-запросов.
	Port int `yaml:"port" env-default:"8080"`
	// Timeout максимальное время ожидания запроса.
	Timeout time.Duration `yaml:"timeout"`
	UseSSL  bool          `yaml:"use_ssl" env-default:"false"`
}

// MustLoad загружает конфиг из yaml-файла с описанием.
// Путь к файлу указывается во флаге запуска config или переменной окружения CONFIG_PATH.
func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not exist " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config is incorrect " + err.Error())
	}
	return &cfg
}

// priority: flag > env > default
func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
