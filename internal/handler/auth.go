package handler

import (
	"net/http"
	"tds_server/internal/config"
	"tds_server/internal/repository"
	"tds_server/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoginRedirect(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authURL := service.BuildAuthURL(cfg)
		c.Redirect(http.StatusFound, authURL)
	}
}

func LoginCallback(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")

		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		}
		teslaTokenRepo, err := service.ExchangeCode(cfg, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// userId := c.Query("state")
		// 保存 token。Save the token.
		saveErr := tokenRepo.Save(uuid.New(), teslaTokenRepo.AccessToken, teslaTokenRepo.RefreshToken, time.Duration(teslaTokenRepo.ExpiresIn))
		if saveErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": saveErr.Error()})
			return
		}

		c.JSON(http.StatusOK, teslaTokenRepo)
	}
}
