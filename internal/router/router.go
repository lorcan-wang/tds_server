package router

import (
	"net/http"
	"path/filepath"
	"runtime"
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

	publicKeyFile := publicKeyFilePath()
	r.StaticFile("/.well-known/appspecific/com.tesla.3p.public-key.pem", publicKeyFile)

	auth := r.Group("/api")
	{
		auth.GET("/login", handler.LoginRedirect(cfg))
		auth.GET("/login/callback", handler.LoginCallback(cfg, tokenRepo))
		auth.GET("/list", handler.GetList(cfg, tokenRepo))
	}
	return r
}

func publicKeyFilePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("public", ".well-known", "appspecific", "com.tesla.3p.public-key.pem")
	}

	baseDir := filepath.Join(filepath.Dir(filename), "..", "..")
	return filepath.Join(baseDir, "public", ".well-known", "appspecific", "com.tesla.3p.public-key.pem")
}
