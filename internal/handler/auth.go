package handler

import (
	"net/http"
	"tds_server/internal/config"
	"tds_server/internal/service"

	"github.com/gin-gonic/gin"
)

func LoginRedirect(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authURL := service.BuildAuthURL(cfg)
		c.Redirect(http.StatusFound, authURL)
	}
}

func LoginCallback(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		}
		token, err := service.ExchangeCode(cfg, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, token)
	}
}
