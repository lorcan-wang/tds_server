package handler

import (
	"encoding/json"
	"fmt"
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

// VehicleListResponse mirrors the Tesla GET /api/1/vehicles payload.
type VehicleListResponse struct {
	// Response holds the fleet vehicle entries in the current page.
	Response []VehicleSummary `json:"response"`
	// Pagination describes the current paging cursor returned by Tesla.
	Pagination PaginationMeta `json:"pagination"`
	// Count is the number of vehicles returned in this payload.
	Count int `json:"count"`
}

// VehicleSummary contains high-level metadata for a single vehicle.
type VehicleSummary struct {
	// ID is the Tesla Fleet unique identifier for the vehicle.
	ID int64 `json:"id"`
	// VehicleID is the vehicle identifier used by mobile access APIs.
	VehicleID int64 `json:"vehicle_id"`
	// VIN is the standard Vehicle Identification Number.
	VIN string `json:"vin"`
	// Color represents the configured exterior color if present.
	Color *string `json:"color"`
	// AccessType indicates the level of access (e.g. OWNER, DRIVER).
	AccessType string `json:"access_type"`
	// DisplayName is the user-defined label for the vehicle.
	DisplayName string `json:"display_name"`
	// OptionCodes is a comma-separated list of factory option codes.
	OptionCodes string `json:"option_codes"`
	// GranularAccess describes fine-grained access controls for the vehicle.
	GranularAccess VehicleGranularAccess `json:"granular_access"`
	// Tokens holds command tokens required by legacy vehicle APIs.
	Tokens []string `json:"tokens"`
	// State reflects the latest connectivity status (e.g. online, asleep).
	State string `json:"state"`
	// InService reports whether the vehicle is currently under service.
	InService bool `json:"in_service"`
	// IDS is the string representation of the vehicle id.
	IDS string `json:"id_s"`
	// CalendarEnabled indicates if calendar sync is available.
	CalendarEnabled bool `json:"calendar_enabled"`
	// APIVersion is the firmware API version exposed via the Fleet API.
	APIVersion *int `json:"api_version"`
	// BackseatToken is populated for backseat account invitations.
	BackseatToken *string `json:"backseat_token"`
	// BackseatTokenUpdatedAt is the last update timestamp for BackseatToken.
	BackseatTokenUpdatedAt *string `json:"backseat_token_updated_at"`
}

// VehicleGranularAccess wraps finer grained controls for shared vehicles.
type VehicleGranularAccess struct {
	// HidePrivate hides private data (for example location) from partner apps.
	HidePrivate bool `json:"hide_private"`
}

// PaginationMeta tracks paging cursors returned by Tesla APIs.
type PaginationMeta struct {
	// Previous is the cursor for the previous page when available.
	Previous *string `json:"previous"`
	// Next is the cursor for the next page when available.
	Next *string `json:"next"`
	// Current is the index of the current page, starting from 1.
	Current int `json:"current"`
	// PerPage is the page size used for the query.
	PerPage int `json:"per_page"`
	// Count is the number of entries returned in the current page.
	Count int `json:"count"`
	// Pages is the total number of pages available for the query.
	Pages int `json:"pages"`
}

// ListVehicles proxies Tesla GET /api/1/vehicles and documents the response payload.
func ListVehicles(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	proxy := newTeslaProxy(cfg, tokenRepo)
	return func(c *gin.Context) {
		query := buildVehicleListQuery(c)
		var payload VehicleListResponse
		status, err := proxy.JSON(c, http.MethodGet, apiSegments("vehicles"), query, nil, nil, &payload)
		if err != nil {
			respondWithError(c, status, err)
			return
		}
		c.JSON(status, payload)
	}
}

// GetVehicle proxies Tesla GET /api/1/vehicles/{vehicle_tag} returning a single vehicle record.
func GetVehicle(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	proxy := newTeslaProxy(cfg, tokenRepo)
	return func(c *gin.Context) {
		var payload VehicleResponse
		status, err := proxy.JSON(c, http.MethodGet, apiSegments("vehicles", ":vehicle_tag"), nil, nil, nil, &payload)
		if err != nil {
			respondWithError(c, status, err)
			return
		}
		c.JSON(status, payload)
	}
}

// VehicleResponse wraps a single vehicle summary payload.
type VehicleResponse struct {
	Response VehicleSummary `json:"response"`
}

// VehicleDataResponse mirrors Tesla GET /api/1/vehicles/{vehicle_tag}/vehicle_data payload.
type VehicleDataResponse struct {
	Response VehicleData `json:"response"`
}

// VehicleData aggregates summary information alongside detailed vehicle states.
type VehicleData struct {
	VehicleSummary
	// UserID is the Tesla account identifier that owns the vehicle.
	UserID int64 `json:"user_id"`
	// UserVehicleBoundAt indicates when the user gained access to this vehicle.
	UserVehicleBoundAt *string `json:"user_vehicle_bound_at"`
	// ChargeState contains live charging information (voltage, current, limits).
	ChargeState map[string]any `json:"charge_state"`
	// ClimateState contains HVAC status such as temperatures and fan levels.
	ClimateState map[string]any `json:"climate_state"`
	// ClosuresState tracks door, window and trunk statuses.
	ClosuresState map[string]any `json:"closures_state"`
	// DriveState contains vehicle movement, heading and speed information.
	DriveState map[string]any `json:"drive_state"`
	// GUISettings reflects unit preferences for distance, time and temperature.
	GUISettings map[string]any `json:"gui_settings"`
	// LocationData holds precise latitude and longitude coordinates.
	LocationData map[string]any `json:"location_data"`
	// ChargeScheduleData reports scheduled charging configuration.
	ChargeScheduleData map[string]any `json:"charge_schedule_data"`
	// PreconditioningScheduleData reports scheduled HVAC preconditioning.
	PreconditioningScheduleData map[string]any `json:"preconditioning_schedule_data"`
	// VehicleConfig lists hardware configuration such as trim, wheels and software packages.
	VehicleConfig map[string]any `json:"vehicle_config"`
	// VehicleState contains alarms, sentry mode, odometer and window/door sensors.
	VehicleState map[string]any `json:"vehicle_state"`
	// VehicleDataCombo is populated when requesting aggregated state bundles.
	VehicleDataCombo map[string]any `json:"vehicle_data_combo"`
}

// VehicleDriverListResponse mirrors Tesla GET /api/1/vehicles/{vehicle_tag}/drivers payload.
type VehicleDriverListResponse struct {
	Response []VehicleDriver `json:"response"`
	Count    int             `json:"count"`
}

// VehicleDriver captures each permitted driver entry.
type VehicleDriver struct {
	MyTeslaUniqueID int64                 `json:"my_tesla_unique_id"`
	UserID          int64                 `json:"user_id"`
	UserIDS         string                `json:"user_id_s"`
	VaultUUID       string                `json:"vault_uuid"`
	DriverFirstName string                `json:"driver_first_name"`
	DriverLastName  string                `json:"driver_last_name"`
	GranularAccess  VehicleGranularAccess `json:"granular_access"`
	ActivePubKeys   []string              `json:"active_pubkeys"`
	PublicKey       string                `json:"public_key"`
}

// GetVehicleData proxies Tesla GET /api/1/vehicles/{vehicle_tag}/vehicle_data for real-time state.
func GetVehicleData(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	proxy := newTeslaProxy(cfg, tokenRepo)
	return func(c *gin.Context) {
		query := buildVehicleDataQuery(c)
		var payload VehicleDataResponse
		status, err := proxy.JSON(c, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "vehicle_data"), query, nil, nil, &payload)
		if err != nil {
			respondWithError(c, status, err)
			return
		}
		c.JSON(status, payload)
	}
}

// GetVehicleDrivers proxies Tesla GET /api/1/vehicles/{vehicle_tag}/drivers to list authorized drivers.
func GetVehicleDrivers(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	proxy := newTeslaProxy(cfg, tokenRepo)
	return func(c *gin.Context) {
		var payload VehicleDriverListResponse
		status, err := proxy.JSON(c, http.MethodGet, apiSegments("vehicles", ":vehicle_tag", "drivers"), nil, nil, nil, &payload)
		if err != nil {
			respondWithError(c, status, err)
			return
		}
		c.JSON(status, payload)
	}
}

type teslaProxy struct {
	cfg       *config.Config
	tokenRepo *repository.TokenRepo
}

func newTeslaProxy(cfg *config.Config, tokenRepo *repository.TokenRepo) *teslaProxy {
	return &teslaProxy{
		cfg:       cfg,
		tokenRepo: tokenRepo,
	}
}

// JSON performs a Tesla API request against the provided path segments and decodes the JSON payload.
func (p *teslaProxy) JSON(
	c *gin.Context,
	method string,
	pathSegments []string,
	query url.Values,
	body []byte,
	headers map[string]string,
	dest any,
) (int, error) {
	path, err := resolveTeslaPath(c, pathSegments...)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resp, status, err := p.do(c, method, path, query, body, headers)
	if err != nil {
		return status, err
	}

	if err := json.Unmarshal(resp.Body(), dest); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("decode Tesla response: %w", err)
	}

	return resp.StatusCode(), nil
}

func (p *teslaProxy) do(
	c *gin.Context,
	method string,
	path string,
	query url.Values,
	body []byte,
	headers map[string]string,
) (*resty.Response, int, error) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		return nil, http.StatusUnauthorized, fmt.Errorf("user is not authenticated")
	}

	token, err := p.tokenRepo.GetByUserID(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if token == nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("user token not found")
	}

	token, err = ensureValidToken(p.cfg, p.tokenRepo, userID, token)
	if err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("token refresh failed: %w", err)
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

	requestURL := buildTeslaURL(p.cfg.TeslaAPIURL, path)

	makeRequest := func(accessToken string) (*resty.Response, error) {
		req := client.R()
		req.SetHeader("Authorization", "Bearer "+accessToken)
		for k, v := range headerValues {
			req.SetHeader(k, v)
		}
		if sanitizedQuery != nil {
			req.SetQueryString(sanitizedQuery.Encode())
		}
		if len(body) > 0 {
			req.SetBody(body)
		}
		return req.Execute(strings.ToUpper(method), requestURL)
	}

	resp, err := makeRequest(token.AccessToken)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		token, err = refreshUserToken(p.cfg, p.tokenRepo, userID, token)
		if err != nil {
			return nil, http.StatusUnauthorized, fmt.Errorf("token refresh failed: %w", err)
		}

		resp, err = makeRequest(token.AccessToken)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if resp.StatusCode() == http.StatusUnauthorized {
			return nil, http.StatusUnauthorized, fmt.Errorf("unauthorized after token refresh")
		}
	}

	if resp.StatusCode() >= http.StatusBadRequest {
		return nil, resp.StatusCode(), fmt.Errorf(resp.String())
	}

	return resp, resp.StatusCode(), nil
}

func respondWithError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

func buildVehicleListQuery(c *gin.Context) url.Values {
	query := url.Values{}
	if page := strings.TrimSpace(c.Query("page")); page != "" {
		query.Set("page", page)
	}
	if perPage := strings.TrimSpace(c.Query("per_page")); perPage != "" {
		query.Set("per_page", perPage)
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func buildVehicleDataQuery(c *gin.Context) url.Values {
	query := url.Values{}
	if endpoints := strings.TrimSpace(c.Query("endpoints")); endpoints != "" {
		query.Set("endpoints", endpoints)
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func apiSegments(segments ...string) []string {
	base := []string{"api", "1"}
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
