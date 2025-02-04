package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SemenShakhray/url-shortener/internal/config"
	"github.com/SemenShakhray/url-shortener/internal/handlers"
	"github.com/SemenShakhray/url-shortener/internal/router"
	"github.com/SemenShakhray/url-shortener/internal/service"
	"github.com/SemenShakhray/url-shortener/internal/storage/postgres"
	"github.com/SemenShakhray/url-shortener/pkg/logger"
)

func main() {

	//TODO init config: cleanenv
	config := config.MustLoad()

	//TODO init logger: slog
	log := logger.SetupLogger(config.Env)
	log.Info("starting url shortener", slog.String("env", config.Env))
	log.Debug("debug massages are enabled")

	//TODO init storage: postgres
	db, err := postgres.Connect(config)
	if err != nil {
		log.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("connected to storage", slog.String("database:", config.DB.User))

	store := postgres.NewStore(db)

	//TODO init service
	service := service.NewService(store)

	//TODO init handlers
	handler := handlers.NewHandler(log.Log, service)

	//TODO init router: gin
	r := router.NewRouter(handler, config)

	//TODO run server
	log.Info("starting server", slog.String("host", config.Server.Host), slog.String("port", config.Server.Port))

	srv := &http.Server{
		Addr:         net.JoinHostPort(config.Server.Host, config.Server.Port),
		Handler:      r,
		ReadTimeout:  config.Server.Timeout,
		WriteTimeout: config.Server.Timeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Error("failed to start server", slog.String("error", err.Error()))
		}
		log.Error("server stoped")
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint

	log.Info(fmt.Sprintf("Received signal: %v", sigint))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("HTTP server shutdown", slog.String("error", err.Error()))
	}
}
