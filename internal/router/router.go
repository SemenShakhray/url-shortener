package router

import (
	"github.com/SemenShakhray/url-shortener/internal/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(h handlers.Handler) *gin.Engine {
	c := gin.Default()

	c.POST("/url", h.SaveURL)

	return c
}
