package handlers

import (
	"log/slog"

	"github.com/SemenShakhray/url-shortener/internal/service"
	"github.com/SemenShakhray/url-shortener/pkg/api/response"
	"github.com/SemenShakhray/url-shortener/pkg/logger"
)

type Handler struct {
	Log  logger.Logger
	Serv service.Servicer
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

func NewHandler(log *slog.Logger, serv service.Servicer) Handler {
	return Handler{
		Serv: serv,
		Log: logger.Logger{
			Log: log,
		},
	}
}
