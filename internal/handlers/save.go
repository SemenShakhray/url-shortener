package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/SemenShakhray/url-shortener/internal/storage"
	"github.com/SemenShakhray/url-shortener/pkg/api/response"
	"github.com/SemenShakhray/url-shortener/pkg/random"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	//TODO: move to config id needed
	aliasLength = 6
)

func (h *Handler) SaveURL(c *gin.Context) {
	op := "handlers.SaveURL"

	resp := Response{}

	log := h.Log.With(
		slog.String("op", op),
	)

	var req Request

	err := c.ShouldBindJSON(&req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			c.JSON(http.StatusBadRequest, resp.Err("empty request"))

			return
		}

		log.Error("failed to decode request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, resp.Err("failed to decode request body"))

		return
	}
	log.Info("request body decoded", slog.Any("request", req))

	err = validator.New().Struct(req)
	if err != nil {
		validateErr := err.(validator.ValidationErrors)

		log.Error("invalid request", slog.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, response.ValidationError(validateErr, resp.Response))

		return
	}

	alias := req.Alias
	if alias == "" {
		alias, err = random.NewRandomString(aliasLength)
		if err != nil {
			log.Error("failed create alias", slog.String("error", err.Error()))

			c.JSON(http.StatusInternalServerError, resp.Err("failed create alias"))

			return
		}
	}

	ctx := c.Request.Context()

	err = h.Serv.SaveURL(ctx, req.URL, alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLOrAliasExists) {
			log.Info("url or alias already exists", slog.String("url", req.URL), slog.String("alias", alias))

			c.JSON(http.StatusBadRequest, resp.Err("url or alias already exists"))

			return
		}
		log.Error("failed add url", slog.String("error", err.Error()))

		c.JSON(http.StatusInternalServerError, resp.Err(err.Error()))

		return
	}
	log.Info("url added", slog.String("url", req.URL), slog.String("alias", alias))

	c.JSON(http.StatusOK, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
