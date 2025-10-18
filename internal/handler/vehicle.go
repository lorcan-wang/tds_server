package handler

import (
	"net/http"
	"net/url"
	"time"

	"tds_server/internal/config"
	"tds_server/internal/model"
	"tds_server/internal/repository"
	"tds_server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const teslaUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"

// GetList returns the vehicles bound to the specified user. GetList 返回指定用户绑定的车辆列表信息。
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
		if token, err = ensureValidToken(cfg, tokenRepo, userID, token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token refresh failed: " + err.Error()})
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
			// Final safety net: if Tesla still reports 401, refresh again to avoid stale tokens.
			// 最终兜底：若 Tesla 仍返回 401，再次刷新以防止使用过期凭证。
			token, err = refreshUserToken(cfg, tokenRepo, userID, token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token refresh failed: " + err.Error()})
				return
			}

			client.SetHeader("Authorization", "Bearer "+token.AccessToken)
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

// GetVehicle retrieves vehicle details for the given tag. GetVehicle 根据车辆标识获取车辆详细信息。
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
		if token, err = ensureValidToken(cfg, tokenRepo, userID, token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token refresh failed: " + err.Error()})
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
			token, err = refreshUserToken(cfg, tokenRepo, userID, token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token refresh failed: " + err.Error()})
				return
			}

			client.SetHeader("Authorization", "Bearer "+token.AccessToken)
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

// GetVehicleData retrieves vehicle_data from the Tesla Fleet API while passing through query parameters (excluding user_id). GetVehicleData 调用 Tesla Fleet API 的 vehicle_data 接口，并透传查询参数（除 user_id 外）。
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
		if token, err = ensureValidToken(cfg, tokenRepo, userID, token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token refresh failed: " + err.Error()})
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
			token, err = refreshUserToken(cfg, tokenRepo, userID, token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token refresh failed: " + err.Error()})
				return
			}

			client.SetHeader("Authorization", "Bearer "+token.AccessToken)
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

// ensureValidToken proactively renews tokens that are about to expire (default 5 minutes window).
// ensureValidToken 会在 token 剩余不足 5 分钟时主动刷新，避免后续请求命中 401。
func ensureValidToken(cfg *config.Config, tokenRepo *repository.TokenRepo, userID uuid.UUID, token *model.UserToken) (*model.UserToken, error) {
	if !token.IsExpired(5 * time.Minute) {
		return token, nil
	}
	return refreshUserToken(cfg, tokenRepo, userID, token)
}

func refreshUserToken(cfg *config.Config, tokenRepo *repository.TokenRepo, userID uuid.UUID, token *model.UserToken) (*model.UserToken, error) {
	// Caller-side 401 guard relies on this function to be idempotent.
	// 调用方的 401 回退依赖此函数幂等刷新，因此这里无需额外判断。
	refreshed, err := service.RefreshToken(cfg, token.RefreshToken)
	if err != nil {
		return nil, err
	}

	newRefresh := refreshed.RefreshToken
	if newRefresh == "" {
		newRefresh = token.RefreshToken
	}

	expiresIn := time.Duration(refreshed.ExpiresIn)
	if err := tokenRepo.Save(userID, refreshed.AccessToken, newRefresh, expiresIn); err != nil {
		return nil, err
	}

	token.AccessToken = refreshed.AccessToken
	token.RefreshToken = newRefresh
	token.ExpiresAt = time.Now().Add(expiresIn * time.Second)
	return token, nil
}
