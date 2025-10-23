package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"tds_server/internal/config"
	"tds_server/internal/repository"
	"tds_server/internal/service"
	"time"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func LoginRedirect(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		state := c.Query("state")
		authURL := service.BuildAuthURL(cfg, state)
		c.Redirect(http.StatusFound, authURL)
	}
}

func LoginCallback(cfg *config.Config, tokenRepo *repository.TokenRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code is required"})
			return
		}

		userID, err := resolveUserID(c.Query("state"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		teslaTokenRepo, err := service.ExchangeCode(cfg, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if saveErr := tokenRepo.Save(userID, teslaTokenRepo.AccessToken, teslaTokenRepo.RefreshToken, time.Duration(teslaTokenRepo.ExpiresIn)); saveErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": saveErr.Error()})
			return
		}

		jwtToken, err := buildJWT(cfg, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := loginCallbackResponse{
			UserID: userID.String(),
			JWT: loginJWT{
				Token:     jwtToken,
				ExpiresIn: int(cfg.JWT.Expiration.Seconds()),
				Issuer:    cfg.JWT.Issuer,
			},
			TeslaToken: teslaTokenRepo,
		}
		fmt.Printf("JWT Token Response: %+v\n", response)
		if prefersJSON(c) {
			c.JSON(http.StatusOK, response)
			return
		}

		if err := renderLoginCallbackHTML(c, response); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

func resolveUserID(state string) (uuid.UUID, error) {
	if state == "" {
		return uuid.New(), nil
	}

	if id, err := uuid.Parse(state); err == nil {
		return id, nil
	}

	// When state is used purely for CSRF/metadata, fall back to creating a new user ID.
	return uuid.New(), nil
}

func buildJWT(cfg *config.Config, userID uuid.UUID) (string, error) {
	if cfg.JWT.Secret == "" {
		return "", errors.New("jwt secret is not configured")
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    cfg.JWT.Issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(cfg.JWT.Expiration)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

type loginCallbackResponse struct {
	UserID     string                      `json:"user_id"`
	JWT        loginJWT                    `json:"jwt"`
	TeslaToken *service.TeslaTokenResponse `json:"tesla_token"`
}

type loginJWT struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	Issuer    string `json:"issuer"`
}

func prefersJSON(c *gin.Context) bool {
	accept := strings.ToLower(c.GetHeader("Accept"))
	if strings.Contains(accept, "application/json") || c.Query("format") == "json" {
		return true
	}
	return false
}

func renderLoginCallbackHTML(c *gin.Context, payload loginCallbackResponse) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(data)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
	<meta charset="utf-8"/>
	<title>登录成功</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; margin: 0; padding: 24px; text-align: center; background: #f7f9fb; color: #1f2933; }
		h1 { font-size: 20px; margin-bottom: 8px; }
		p { font-size: 14px; margin: 0; }
	</style>
</head>
<body>
	<h1>登录成功</h1>
	<p>可以返回应用，页面会在片刻后自动关闭。</p>
	<script>
	var __TDSPayloadB64 = %q;
	(function() {
		try {
			var payload = JSON.parse(atob(__TDSPayloadB64));
			var deepLink = "tdsclient://auth/callback?payload=" + encodeURIComponent(__TDSPayloadB64);
			if (window.ReactNativeWebView && window.ReactNativeWebView.postMessage) {
				window.ReactNativeWebView.postMessage(JSON.stringify(payload));
			} else if (window.opener && window.opener.postMessage) {
				window.opener.postMessage(payload, "*");
			}
			setTimeout(function () {
				try {
					if (window.location && typeof window.location.replace === "function") {
						window.location.replace(deepLink);
					} else {
						window.location.href = deepLink;
					}
				} catch (err) {
					console.error("Failed to redirect to client scheme", err);
				}
			}, 150);
		} catch (err) {
			console.error("Failed to deliver login payload", err);
		}
		setTimeout(function () {
			if (window.close) {
				window.close();
			}
		}, 800);
	})();
	</script>
	<noscript>需要启用 JavaScript 以完成登录，可以手动返回应用。</noscript>
</body>
</html>`, encoded)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	return nil
}
