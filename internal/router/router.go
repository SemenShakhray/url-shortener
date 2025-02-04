package router

import (
	"net/http"

	"github.com/SemenShakhray/url-shortener/internal/config"
	"github.com/SemenShakhray/url-shortener/internal/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter(h handlers.Handler, cfg config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/url/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "alias cannot be empty"})
	})
	r.GET("url/:alias", h.Redirect)

	rAuth := r.Group("/url", gin.BasicAuth(gin.Accounts{cfg.Server.User: cfg.Server.Password}))
	rAuth.POST("", h.SaveURL)
	rAuth.DELETE("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "alias cannot be empty"})
	})
	rAuth.DELETE("/:alias", h.DeleteURL)

	return r
}
