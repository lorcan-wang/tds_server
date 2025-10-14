package handler

import (
	"net/http"
	"net/url"
	"time"

	"tds_server/internal/config"
	"tds_server/internal/repository"
	"tds_server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const teslaUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"

// GetList 返回指定用户绑定的车辆列表信息。
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

		contentType := resp.Header().Get("Content-Type")
		if contentType == "" {
			contentType = "application/json"
		}
		c.Data(http.StatusOK, contentType, resp.Body())
	}
}

// GetVehicle 根据车辆标识获取车辆详细信息。
func GetVehicle(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDParam, ok := c.GetQuery("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		vehicleTag := c.Param("vehicle_tag")
		if vehicleTag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_tag is required"})
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

		detailURL := cfg.TeslaAPIURL + "/api/1/vehicles/" + url.PathEscape(vehicleTag)
		resp, err := client.R().Get(detailURL)
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
			resp, err = client.R().Get(detailURL)
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

		contentType := resp.Header().Get("Content-Type")
		if contentType == "" {
			contentType = "application/json"
		}
		c.Data(http.StatusOK, contentType, resp.Body())
	}
}

// GetVehicleData 调用 Tesla Fleet API 的 vehicle_data 接口，支持透传查询参数（除 user_id 外）。
func GetVehicleData(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDParam, ok := c.GetQuery("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		vehicleTag := c.Param("vehicle_tag")
		if vehicleTag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_tag is required"})
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

		requestURL := cfg.TeslaAPIURL + "/api/1/vehicles/" + url.PathEscape(vehicleTag) + "/vehicle_data"

		query := c.Request.URL.Query()
		query.Del("user_id")

		makeRequest := func() (*resty.Response, error) {
			req := client.R()
			if len(query) > 0 {
				req.SetQueryString(query.Encode())
			}
			return req.Get(requestURL)
		}

		resp, err := makeRequest()
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
			resp, err = makeRequest()
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

		contentType := resp.Header().Get("Content-Type")
		if contentType == "" {
			contentType = "application/json"
		}
		c.Data(http.StatusOK, contentType, resp.Body())
	}
}
