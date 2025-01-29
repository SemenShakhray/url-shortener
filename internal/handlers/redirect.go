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

	h.Log.Log.With(
		slog.String("op", op),
	)

	ctx := c.Request.Context()
	alias := c.Param("alias")

	url, err := h.Serv.GetURL(ctx, alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			h.Log.Error("url not found", slog.String("alias", alias))
			c.JSON(http.StatusBadRequest, resp.Err("url not found"))
			return
		}

		h.Log.Error("failed to get url", slog.String("alias", alias))
		c.JSON(http.StatusInternalServerError, resp.Err("internal error"))
		return
	}

	h.Log.Info("got URL", slog.String("url", url))

	c.Redirect(http.StatusFound, url)
}
