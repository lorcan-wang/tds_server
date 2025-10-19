package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"tds_server/internal/config"
	"tds_server/internal/middleware"
	"tds_server/internal/repository"
	"tds_server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// VehicleCommand handles Tesla vehicle command requests via POST. VehicleCommand 统一处理 Tesla 车辆指令调用，所有指令均通过 POST 方式触发。
func VehicleCommand(cfg *config.Config, tokenRepo *repository.TokenRepo, commandSvc *service.VehicleCommandService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
			return
		}

		commandPath := strings.Trim(c.Param("command_path"), "/")
		if commandPath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "command_path is required"})
			return
		}

		vehicleTag := c.Param("vehicle_tag")
		if vehicleTag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_tag is required"})
			return
		}

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

		token, err = ensureValidToken(cfg, tokenRepo, userID, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token refresh failed: " + err.Error()})
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
			return
		}

		var commandResult *service.CommandResult
		if commandSvc != nil {
			commandName := strings.Split(commandPath, "/")[0]
			commandResult, err = commandSvc.Execute(c.Request.Context(), vehicleTag, commandName, bodyBytes, token.AccessToken)
			switch {
			case err == nil && commandResult != nil:
				c.Data(commandResult.Status, commandResult.ContentType, commandResult.Body)
				return
			case errors.Is(err, service.ErrVehicleCommandUseREST):
				// fall back to REST handling below. 在下方回退到 REST 处理。
			case err != nil:
				var cmdErr *service.CommandError
				if errors.As(err, &cmdErr) {
					if len(cmdErr.Body) > 0 {
						c.Data(cmdErr.Status, "application/json", cmdErr.Body)
					} else {
						c.JSON(cmdErr.Status, gin.H{"error": cmdErr.Error()})
					}
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		requestURL := buildVehicleCommandURL(cfg.TeslaAPIURL, vehicleTag, commandPath)

		query := c.Request.URL.Query()
		query.Del("user_id")

		makeRequest := func(accessToken string) (*resty.Response, error) {
			client := resty.New()
			client.SetContentLength(true)
			client.SetHeader("User-Agent", teslaUserAgent)

			req := client.R()
			req.SetHeader("Authorization", "Bearer "+accessToken)

			if accept := c.GetHeader("Accept"); accept != "" {
				req.SetHeader("Accept", accept)
			}
			if contentType := c.GetHeader("Content-Type"); contentType != "" {
				req.SetHeader("Content-Type", contentType)
			} else if len(bodyBytes) > 0 {
				req.SetHeader("Content-Type", "application/json")
			}
			if len(bodyBytes) > 0 {
				req.SetBody(bodyBytes)
			}
			if len(query) > 0 {
				req.SetQueryString(query.Encode())
			}
			return req.Post(requestURL)
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
}

func buildVehicleCommandURL(baseURL, vehicleTag, commandPath string) string {
	base := strings.TrimRight(baseURL, "/")
	escapedVehicleTag := url.PathEscape(vehicleTag)

	segments := strings.Split(commandPath, "/")
	escapedSegments := make([]string, 0, len(segments))
	for _, seg := range segments {
		if seg == "" {
			continue
		}
		escapedSegments = append(escapedSegments, url.PathEscape(seg))
	}

	command := strings.Join(escapedSegments, "/")
	return fmt.Sprintf("%s/api/1/vehicles/%s/command/%s", base, escapedVehicleTag, command)
}
