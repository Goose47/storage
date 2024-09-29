package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Mode     string
	Host     string
	Port     string
	Secret   string
	AuthHost string
	AuthPort string

	DB *DBConfig
	FS *FSConfig
}

type DBConfig struct {
	Url    string
	DBName string
	DBColl string
}

type FSConfig struct {
	Base string
}

const configPath = ".env"

// MustLoad loads config from .env file and panics on error
func MustLoad() *Config {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("config file does not exist: %s", configPath))
	}

	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("failed to load config file: %w", err))
	}

	cfg := Config{
		Mode:     mustRetrieve("APP_MODE"),
		Host:     mustRetrieve("APP_HOST"),
		Port:     mustRetrieve("APP_PORT"),
		Secret:   mustRetrieve("APP_SECRET"),
		AuthHost: mustRetrieve("APP_AUTH_HOST"),
		AuthPort: mustRetrieve("APP_AUTH_PORT"),
	}

	cfg.DB = &DBConfig{
		Url:    mustRetrieve("DB_URL"),
		DBName: mustRetrieve("DB_NAME"),
		DBColl: mustRetrieve("DB_COLL"),
	}

	cfg.FS = &FSConfig{
		Base: mustRetrieve("STORAGE_PATH"),
	}

	return &cfg
}

// mustRetrieve retrieves key from environment variables and panics if value is empty
func mustRetrieve(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Errorf("%s is not present in env", key))
	}
	return val
}
