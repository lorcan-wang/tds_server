package handler

import (
	"net/http"
	"tds_server/internal/config"
	"tds_server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

func GetList(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.GetQuery("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		token, err := tokenRepo.GetByUserID(uuid.MustParse(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		client := resty.New()
		client.SetHeader("Authorization", "Bearer "+token.AccessToken)
		client.SetContentLength(true)
		client.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
		resp, err := client.R().Get(cfg.TeslaAPIURL + "/api/1/vehicles")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": resp.String()})
	}
}
