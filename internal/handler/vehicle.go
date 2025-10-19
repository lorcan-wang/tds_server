package handler

import (
	"bytes"
	"fmt"
	"io"
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

// Vehicle inventory endpoints --------------------------------------------------------------------
// 车辆清单相关端点：负责车辆基础信息与状态拉取。

func GetList(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles")...)
}

func GetVehicle(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag")...)
}

func GetVehicleData(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "vehicle_data")...)
}

func GetVehicleStates(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "states")...)
}

func GetVehicleDataRequest(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "vehicle_data_request", ":state")...)
}

func WakeUpVehicle(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodPost, apiSegments("vehicles", ":vehicle_tag", "wake_up")...)
}

// Driver and sharing endpoints -------------------------------------------------------------------
// 驾驶员与车辆共享端点：管理共享邀请与驾驶授权。

func GetVehicleDrivers(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "drivers")...)
}

func DeleteVehicleDriver(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodDelete, apiSegments("vehicles", ":vehicle_tag", "drivers")...)
}

func GetVehicleShareInvites(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "invitations")...)
}

func CreateVehicleShareInvite(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodPost, apiSegments("vehicles", ":vehicle_tag", "invitations")...)
}

func RevokeVehicleShareInvite(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodPost, apiSegments("vehicles", ":vehicle_tag", "invitations", ":invitation_id", "revoke")...)
}

func RedeemShareInvite(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodPost, apiSegments("invitations", "redeem")...)
}

// Vehicle capability endpoints -------------------------------------------------------------------
// 能力查询端点：涵盖附近充电站、固件信息等扩展能力。

func GetVehicleMobileEnabled(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "mobile_enabled")...)
}

func GetNearbyChargingSites(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "nearby_charging_sites")...)
}

func GetRecentAlerts(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "recent_alerts")...)
}

func GetReleaseNotes(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "release_notes")...)
}

func GetServiceData(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "service_data")...)
}

// Fleet telemetry endpoints ----------------------------------------------------------------------
// 车队遥测端点：配置/查询第三方遥测服务。

func CreateFleetTelemetryConfig(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodPost, apiSegments("vehicles", "fleet_telemetry_config")...)
}

func GetFleetTelemetryConfig(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "fleet_telemetry_config")...)
}

func DeleteFleetTelemetryConfig(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodDelete, apiSegments("vehicles", ":vehicle_tag", "fleet_telemetry_config")...)
}

func PostFleetTelemetryConfigJWS(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodPost, apiSegments("vehicles", "fleet_telemetry_config_jws")...)
}

func GetFleetTelemetryErrors(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "fleet_telemetry_errors")...)
}

// Subscription endpoints -------------------------------------------------------------------------
// 推送订阅端点：维护车辆/设备推送通知订阅列表。

func GetVehicleSubscriptions(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("vehicle_subscriptions")...)
}

func SetVehicleSubscriptions(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodPost, apiSegments("vehicle_subscriptions")...)
}

func GetSubscriptions(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiSegments("subscriptions")...)
}

func SetSubscriptions(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodPost, apiSegments("subscriptions")...)
}

func GetEligibleSubscriptions(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiDXSegments("vehicles", "subscriptions", "eligibility")...)
}

func GetEligibleUpgrades(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiDXSegments("vehicles", "upgrades", "eligibility")...)
}

// Misc endpoints ---------------------------------------------------------------------------------
// 其他车辆相关端点：包含状态汇总、签名指令等。

func PostFleetStatus(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodPost, apiSegments("vehicles", "fleet_status")...)
}

func GetVehicleOptions(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiDXSegments("vehicles", "options")...)
}

func GetWarrantyDetails(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodGet, apiDXSegments("warranty", "details")...)
}

func SendSignedVehicleCommand(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return vehicleEndpointHandler(cfg, tokenRepo, http.MethodPost, apiSegments("vehicles", ":vehicle_tag", "signed_command")...)
}

// Shared proxy helpers ---------------------------------------------------------------------------
// 通用代理工具：统一处理路径解析、令牌刷新与请求透传。

func vehicleEndpointHandler(cfg *config.Config, tokenRepo *repository.TokenRepo, method string, segments ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path, err := resolveTeslaPath(c, segments...)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		body, err := readRequestBody(c, method)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		proxyTeslaRequest(c, cfg, tokenRepo, method, path, c.Request.URL.Query(), body, nil)
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
	if headers != nil {
		for k, v := range headers {
			if v != "" {
				headerValues[k] = v
			}
		}
	}

	if accept := c.GetHeader("Accept"); accept != "" {
		if _, exists := headerValues["Accept"]; !exists {
			headerValues["Accept"] = accept
		}
	}

	if len(body) > 0 && !strings.EqualFold(method, http.MethodGet) {
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

func apiSegments(segments ...string) []string {
	base := []string{"api", "1"}
	return append(base, segments...)
}

func apiDXSegments(segments ...string) []string {
	base := []string{"api", "1", "dx"}
	return append(base, segments...)
}

func resolveTeslaPath(c *gin.Context, segments ...string) (string, error) {
	resolved := make([]string, 0, len(segments))
	for _, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			key := strings.TrimPrefix(segment, ":")
			value := c.Param(key)
			if value == "" {
				return "", fmt.Errorf("%s is required", key)
			}
			resolved = append(resolved, url.PathEscape(value))
			continue
		}
		resolved = append(resolved, segment)
	}
	return "/" + strings.Join(resolved, "/"), nil
}

func readRequestBody(c *gin.Context, method string) ([]byte, error) {
	if strings.EqualFold(method, http.MethodGet) {
		return nil, nil
	}
	if c.Request.Body == nil {
		return nil, nil
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, nil
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

func buildTeslaURL(base, path string) string {
	base = strings.TrimRight(base, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
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
