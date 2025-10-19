package handler

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"tds_server/internal/config"
	"tds_server/internal/middleware"
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
		proxyTeslaRequest(c, cfg, tokenRepo, http.MethodGet, buildVehicleResourcePath(""), c.Request.URL.Query(), nil, nil)
	}
}

// GetVehicle retrieves vehicle details for the given tag. GetVehicle 根据车辆标识获取车辆详细信息。
func GetVehicle(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		vehicleTag := c.Param("vehicle_tag")
		if vehicleTag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_tag is required"})
			return
		}

		path := buildVehicleResourcePath(vehicleTag)
		proxyTeslaRequest(c, cfg, tokenRepo, http.MethodGet, path, c.Request.URL.Query(), nil, nil)
	}
}

// GetVehicleData retrieves vehicle_data from the Tesla Fleet API while passing through query parameters (excluding user_id). GetVehicleData 调用 Tesla Fleet API 的 vehicle_data 接口，并透传查询参数（除 user_id 外）。
func GetVehicleData(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		vehicleTag := c.Param("vehicle_tag")
		if vehicleTag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_tag is required"})
			return
		}

		path := buildVehicleResourcePath(vehicleTag, "vehicle_data")
		proxyTeslaRequest(c, cfg, tokenRepo, http.MethodGet, path, c.Request.URL.Query(), nil, nil)
	}
}

// GetVehicleStates 获取车辆状态列表示例，映射 Tesla Fleet API `/api/1/vehicles/{vehicle_id}/states`。
func GetVehicleStates(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		vehicleTag := c.Param("vehicle_tag")
		if vehicleTag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_tag is required"})
			return
		}

		path := buildVehicleResourcePath(vehicleTag, "states")
		proxyTeslaRequest(c, cfg, tokenRepo, http.MethodGet, path, c.Request.URL.Query(), nil, nil)
	}
}

// GetVehicleDataRequest 透传 `vehicle_data_request/{state}` 接口，支持 charge_state、drive_state 等细分数据。
func GetVehicleDataRequest(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		vehicleTag := c.Param("vehicle_tag")
		if vehicleTag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_tag is required"})
			return
		}

		state := c.Param("state")
		if state == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "state is required"})
			return
		}

		path := buildVehicleResourcePath(vehicleTag, "vehicle_data_request", state)
		proxyTeslaRequest(c, cfg, tokenRepo, http.MethodGet, path, c.Request.URL.Query(), nil, nil)
	}
}

// WakeUpVehicle 转发官方 `/api/1/vehicles/{vehicle_id}/wake_up` 接口，确保车辆唤醒。
func WakeUpVehicle(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		vehicleTag := c.Param("vehicle_tag")
		if vehicleTag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_tag is required"})
			return
		}

		path := buildVehicleResourcePath(vehicleTag, "wake_up")
		proxyTeslaRequest(c, cfg, tokenRepo, http.MethodPost, path, c.Request.URL.Query(), nil, nil)
	}
}

func proxyTeslaRequest(
	c *gin.Context,
	cfg *config.Config,
	tokenRepo *repository.TokenRepo,
	method string,
	path string,
	query url.Values,
	body []byte,
	headers map[string]string,
) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not authenticated"})
		return
	}

	token, err := tokenRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if token == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user token not found"})
		return
	}

	if token, err = ensureValidToken(cfg, tokenRepo, userID, token); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token refresh failed: " + err.Error()})
		return
	}

	sanitizedQuery := sanitizeQuery(query)

	client := resty.New()
	client.SetContentLength(true)
	client.SetHeader("User-Agent", teslaUserAgent)

	headerValues := map[string]string{}
	for k, v := range headers {
		if v != "" {
			headerValues[k] = v
		}
	}

	if accept := c.GetHeader("Accept"); accept != "" {
		if _, exists := headerValues["Accept"]; !exists {
			headerValues["Accept"] = accept
		}
	}

	if len(body) > 0 && strings.EqualFold(method, http.MethodPost) {
		if contentType := c.GetHeader("Content-Type"); contentType != "" {
			if _, exists := headerValues["Content-Type"]; !exists {
				headerValues["Content-Type"] = contentType
			}
		} else if _, exists := headerValues["Content-Type"]; !exists {
			headerValues["Content-Type"] = "application/json"
		}
	}

	requestURL := buildTeslaURL(cfg.TeslaAPIURL, path)

	makeRequest := func(accessToken string) (*resty.Response, error) {
		req := client.R()
		req.SetHeader("Authorization", "Bearer "+accessToken)
		for k, v := range headerValues {
			req.SetHeader(k, v)
		}
		if len(sanitizedQuery) > 0 {
			req.SetQueryString(sanitizedQuery.Encode())
		}
		if len(body) > 0 {
			req.SetBody(body)
		}
		return req.Execute(strings.ToUpper(method), requestURL)
	}

	resp, err := makeRequest(token.AccessToken)
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

		resp, err = makeRequest(token.AccessToken)
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
	c.Data(resp.StatusCode(), contentType, resp.Body())
}

func buildTeslaURL(base, path string) string {
	base = strings.TrimRight(base, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}

func buildVehicleResourcePath(vehicleTag string, segments ...string) string {
	parts := []string{"api", "1", "vehicles"}
	if vehicleTag != "" {
		parts = append(parts, url.PathEscape(vehicleTag))
	}
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		parts = append(parts, url.PathEscape(segment))
	}
	return "/" + strings.Join(parts, "/")
}

func sanitizeQuery(values url.Values) url.Values {
	if len(values) == 0 {
		return nil
	}

	filtered := url.Values{}
	for key, val := range values {
		if strings.EqualFold(key, "user_id") {
			continue
		}
		for _, item := range val {
			filtered.Add(key, item)
		}
	}

	if len(filtered) == 0 {
		return nil
	}
	return filtered
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
