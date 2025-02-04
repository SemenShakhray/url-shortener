package handlers

import (
	"log/slog"
	"net/http"

	"github.com/SemenShakhray/url-shortener/pkg/api/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteURL(c *gin.Context) {
	var resp response.Response

	op := "handlers.DeleteURL"

	log := h.Log.With(
		slog.String("op", op),
	)

	alias := c.Param("alias")

	err := h.Serv.DeleteURL(c.Request.Context(), alias)
	if err != nil {
		log.Error("failed to delete url", slog.String("alias", alias), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, resp.Err("failed delete url"))
		return
	}

	c.JSON(http.StatusOK, resp.OK())
}
