package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SemenShakhray/url-shortener/internal/service"
	"github.com/SemenShakhray/url-shortener/internal/storage"
	"github.com/SemenShakhray/url-shortener/pkg/logger"
	"github.com/SemenShakhray/url-shortener/pkg/random"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

const (
	statusOK    = "OK"
	statusError = "Error"

	//TODO: move to config id needed
	aliasLength = 6
)

func (resp *Response) OK() Response {
	return Response{
		Status: statusOK,
	}
}

func (resp *Response) Err(msg string) Response {
	return Response{
		Status: statusError,
		Error:  msg,
	}
}

func NewHandler(log *slog.Logger, serv service.Servicer) Handler {
	return Handler{
		Serv: serv,
		Log: logger.Logger{
			Log: log,
		},
	}
}

func (h *Handler) SaveURL(c *gin.Context) {
	op := "handlers.SaveURL"

	resp := Response{}

	h.Log.Log.With(
		slog.String("op", op),
	)

	var req Request

	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.Log.Error("failed to decode request body", slog.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, resp.Err("failed to decode request body"))
		return
	}
	h.Log.Info("request body decoded", slog.Any("requset", req))

	err = validator.New().Struct(req)
	if err != nil {
		validateErr := err.(validator.ValidationErrors)

		h.Log.Error("invalid request", slog.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, ValidationError(validateErr, resp))
		return
	}

	alias := req.Alias
	if alias == "" {
		alias, err = random.NewRandomString(aliasLength)
		if err != nil {
			h.Log.Error("failed create alias", slog.String("error", err.Error()))
		}
	}

	ctx := c.Request.Context()

	err = h.Serv.SaveURL(ctx, req.URL, alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLOrAliasExists) {
			h.Log.Info("url or alias already exists", slog.String("url", req.URL), slog.String("alias", alias))

			c.JSON(http.StatusBadRequest, resp.Err("url or alias already exists"))
			return
		}
		h.Log.Error("failed add url", slog.String("error", err.Error()))

		c.JSON(http.StatusInternalServerError, resp.Err("failed add url"))
	}

	h.Log.Info("url added", slog.String("url", req.URL), slog.String("alias", alias))

	c.JSON(http.StatusOK, Response{
		Status: resp.OK().Status,
		Alias:  alias,
	})
}

func ValidationError(errs validator.ValidationErrors, resp Response) Response {
	var buff bytes.Buffer

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			buff.WriteString(fmt.Sprintf("|| field %s is required field ||", err.Field()))
		case "url":
			buff.WriteString(fmt.Sprintf("|| field %s is not valid URL ||", err.Field()))
		default:
			buff.WriteString(fmt.Sprintf("|| field %s is not valid ||", err.Field()))
		}
	}
	return resp.Err(buff.String())
}
