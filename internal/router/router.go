package router

import (
	"tds_server/internal/config"
	"tds_server/internal/handler"
	"tds_server/internal/repository"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, tokenRepo *repository.TokenRepo) *gin.Engine {
	r := gin.New()

	auth := r.Group("/api")
	{
		auth.GET("/login", handler.LoginRedirect(cfg))
		auth.GET("/login/callback", handler.LoginCallback(cfg, tokenRepo))
	}
	return r
}
