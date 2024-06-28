package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env         string
	StoragePath string
	HTTPServer
}

type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
	Username    string
	Password    string
}

func MustLoad() *Config {
	const baseError = "error loading .env file"
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("%s: %s", baseError, err.Error())
	// }

	var cfg Config

	// ENV
	cfg.Env = os.Getenv("ENV")
	if cfg.Env == "" {
		log.Fatalf("%s: ENV is required", baseError)
	}
	if cfg.Env != "local" && cfg.Env != "dev" && cfg.Env != "prod" {
		log.Fatalf("%s: invalid ENV value. Must be 'local', 'dev' or 'prod'", baseError)
	}

	// STORAGE_PATH
	cfg.StoragePath = os.Getenv("STORAGE_PATH")
	if cfg.Env == "" {
		log.Fatalf("%s: STORAGE_PATH is required", baseError)
	}

	// SERVER_ADDRESS
	cfg.HTTPServer.Address = os.Getenv("SERVER_ADDRESS")
	if cfg.Env == "" {
		log.Fatalf("%s: SERVER_ADDRESS is required", baseError)
	}

	//SERVER_ADMIN_USERNAME
	cfg.HTTPServer.Username = os.Getenv("SERVER_ADMIN_USERNAME")
	if cfg.Env == "" {
		log.Fatalf("%s: SERVER_ADMIN_USERNAME is required", baseError)
	}

	//SERVER_ADMIN_PASSWORD
	cfg.HTTPServer.Password = os.Getenv("SERVER_ADMIN_PASSWORD")
	if cfg.Env == "" {
		log.Fatalf("%s: SERVER_ADMIN_PASSWORD is required", baseError)
	}

	//SERVER_TIMEOUT
	timeOutString := os.Getenv("SERVER_TIMEOUT")
	if timeOutString == "" {
		log.Fatalf("%s: SERVER_TIMEOUT is required", baseError)
	}
	timeOutInt, err := strconv.Atoi(timeOutString)
	if err != nil {
		log.Fatalf("%s: invalid SERVER_TIMEOUT value. Must be int type", baseError)
	}
	cfg.HTTPServer.Timeout = time.Duration(timeOutInt) * time.Second

	//SERVER_IDLE_TIMEOUT
	idleTimeOutString := os.Getenv("SERVER_IDLE_TIMEOUT")
	if idleTimeOutString == "" {
		log.Fatalf("%s: SERVER_IDLE_TIMEOUT is required", baseError)
	}
	idleTimeOutInt, err := strconv.Atoi(idleTimeOutString)
	if err != nil {
		log.Fatalf("%s: invalid SERVER_IDLE_TIMEOUT value. Must be int type", baseError)
	}
	cfg.HTTPServer.IdleTimeout = time.Duration(idleTimeOutInt) * time.Second

	return &cfg
}
