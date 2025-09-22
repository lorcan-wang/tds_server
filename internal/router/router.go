package router

import (
	"tds_server/internal/config"
	"tds_server/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()

	auth := r.Group("/api")
	{
		auth.GET("/login", handler.LoginRedirect(cfg))
		auth.GET("/login/callback", handler.LoginCallback(cfg))
	}
	return r
}
