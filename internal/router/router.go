package router

import (
	"net/http"
	"path/filepath"
	"runtime"
	"tds_server/internal/config"
	"tds_server/internal/handler"
	"tds_server/internal/repository"
	"tds_server/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, tokenRepo *repository.TokenRepo, partnerSvc *service.PartnerTokenService) *gin.Engine {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ping"})
	})

	if partnerSvc != nil {
		r.Use(func(c *gin.Context) {
			c.Set("partnerTokenService", partnerSvc)
			c.Next()
		})
	}

	publicKeyFile := publicKeyFilePath()
	r.StaticFile("/.well-known/appspecific/com.tesla.3p.public-key.pem", publicKeyFile)

	auth := r.Group("/api")
	{
		auth.GET("/login", handler.LoginRedirect(cfg))
		auth.GET("/login/callback", handler.LoginCallback(cfg, tokenRepo))
		auth.GET("/list", handler.GetList(cfg, tokenRepo))
		auth.GET("/vehicles/:vehicle_tag", handler.GetVehicle(cfg, tokenRepo))
		auth.GET("/vehicles/:vehicle_tag/vehicle_data", handler.GetVehicleData(cfg, tokenRepo))
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
