package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/goloop/env"
)

type Config struct {
	Env       string
	DB        DB
	Server    Server
	Migration Migration
	Clients   ClientsConfig
	AppSecret string
}

type DB struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type Server struct {
	Host        string
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
	User        string
	Password    string
}

type Client struct {
	Address      string
	Timeout      time.Duration
	RetriesCount int
}

type ClientsConfig struct {
	SSO Client
}

type Migration struct {
	Dir  string
	DSN  string
	Name string
}

func MustLoad() Config {
	cfg := Config{}
	err := env.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}
	cfg = Config{
		Env: env.Get("ENV"),

		DB: DB{
			Host: env.Get("DB_HOST"),
			Port: env.Get("DB_PORT"),
			User: env.Get("DB_USER"),
			Pass: env.Get("DB_PASS"),
			Name: env.Get("DB_NAME"),
		},

		Server: Server{
			Host:     env.Get("SERVER_HOST"),
			Port:     env.Get("SERVER_PORT"),
			User:     env.Get("SERVER_USER"),
			Password: env.Get("SERVER_PASSWORD"),
		},
		Migration: Migration{
			Dir:  env.Get("MIGRATION_DIR"),
			DSN:  env.Get("MIGRATION_DSN"),
			Name: env.Get("MIGRATION_NAME"),
		},
		Clients: ClientsConfig{
			SSO: Client{
				Address: env.Get("CLIENT_ADDRESS"),
			},
		},
		AppSecret: env.Get("APP_SECRET"),
	}

	timeoutServer, err := time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	idleTimeoutServer, err := time.ParseDuration(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	cfg.Server.Timeout = timeoutServer
	cfg.Server.IdleTimeout = idleTimeoutServer

	timeoutClient, err := time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}
	retriesCount, err := strconv.Atoi(env.Get("CLIENT_RETRIES_COUNT"))
	if err != nil {
		log.Fatal(err)
	}

	cfg.Clients.SSO.Timeout = timeoutClient
	cfg.Clients.SSO.RetriesCount = retriesCount

	return cfg
}
