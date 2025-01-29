package router

import (
	"net/http"

	"github.com/SemenShakhray/url-shortener/internal/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter(h handlers.Handler) *gin.Engine {
	c := gin.New()

	c.Use(gin.Recovery())

	c.POST("/url", h.SaveURL)
	c.Any("/url/", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alias cannot be empty"})
	})
	c.GET("/url/:alias", h.Redirect)
	c.DELETE("/url/:alias", h.DeleteURL)

	return c
}
