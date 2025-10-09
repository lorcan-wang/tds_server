package router

import (
	"net/http"
	"tds_server/internal/config"
	"tds_server/internal/handler"
	"tds_server/internal/repository"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, tokenRepo *repository.TokenRepo) *gin.Engine {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ping"})
	})
	auth := r.Group("/api")
	{
		auth.GET("/login", handler.LoginRedirect(cfg))
		auth.GET("/login/callback", handler.LoginCallback(cfg, tokenRepo))
		auth.GET("/list", handler.GetList(cfg, tokenRepo))
	}
	return r
}
