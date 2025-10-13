package handler

import (
	"net/http"
	"time"

	"tds_server/internal/config"
	"tds_server/internal/repository"
	"tds_server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const teslaUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"

func GetList(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDParam, ok := c.GetQuery("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		userID, err := uuid.Parse(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is invalid"})
			return
		}

		token, err := tokenRepo.GetByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		client := resty.New()
		client.SetContentLength(true)
		client.SetHeader("User-Agent", teslaUserAgent)
		client.SetHeader("Authorization", "Bearer "+token.AccessToken)

		listURL := cfg.TeslaAPIURL + "/api/1/vehicles"
		resp, err := client.R().Get(listURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if resp.StatusCode() == http.StatusUnauthorized {
			refreshed, refreshErr := service.RefreshToken(cfg, token.RefreshToken)
			if refreshErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token refresh failed: " + refreshErr.Error()})
				return
			}

			newRefresh := refreshed.RefreshToken
			if newRefresh == "" {
				newRefresh = token.RefreshToken
			}

			if saveErr := tokenRepo.Save(userID, refreshed.AccessToken, newRefresh, time.Duration(refreshed.ExpiresIn)); saveErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": saveErr.Error()})
				return
			}

			client.SetHeader("Authorization", "Bearer "+refreshed.AccessToken)
			resp, err = client.R().Get(listURL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if resp.StatusCode() == http.StatusUnauthorized {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized after token refresh"})
				return
			}
		}

		if resp.StatusCode() >= http.StatusBadRequest {
			c.JSON(resp.StatusCode(), gin.H{"error": resp.String()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": resp.String()})
	}
}
