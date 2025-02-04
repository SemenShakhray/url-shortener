package config

import (
	"log"
	"os"
	"time"

	"github.com/goloop/env"
)

type Config struct {
	Env       string
	DB        DB
	Server    Server
	Migration Migration
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
	}

	timeout, err := time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	cfg.Server.Timeout = timeout
	cfg.Server.IdleTimeout = idleTimeout

	return cfg
}
