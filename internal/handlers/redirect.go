package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/SemenShakhray/url-shortener/internal/storage"
	"github.com/SemenShakhray/url-shortener/pkg/api/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Redirect(c *gin.Context) {
	const op = "handlers.Redirect"

	resp := response.Response{}

	log := h.Log.With(
		slog.String("op", op),
	)

	ctx := c.Request.Context()
	alias := c.Param("alias")

	url, err := h.Serv.GetURL(ctx, alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("url not found", slog.String("alias", alias), slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, resp.Err("url not found"))
			return
		}

		log.Error("failed to get url", slog.String("alias", alias), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resp.Err("internal error"))
		return
	}

	log.Debug("got URL", slog.String("url", url))

	c.Redirect(http.StatusFound, url)
}
