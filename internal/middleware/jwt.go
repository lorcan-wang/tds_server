package middleware

import (
	"errors"
	"net/http"
	"strings"
	"tds_server/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const UserIDContextKey = "userID"

// JWTAuth 校验 Authorization 头中的 Bearer JWT，并将用户 UUID 注入 Gin 上下文。
func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if cfg.JWT.Secret == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "jwt secret is not configured"})
			return
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token subject must be a UUID"})
			return
		}

		c.Set(UserIDContextKey, userID)
		c.Next()
	}
}

// UserIDFromContext 从 Gin 上下文读取用户 UUID。
func UserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	if value, ok := c.Get(UserIDContextKey); ok {
		if id, valid := value.(uuid.UUID); valid {
			return id, true
		}
	}
	return uuid.Nil, false
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("authorization header is required")
	}

	if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return "", errors.New("authorization header must be Bearer token")
	}

	token := strings.TrimSpace(header[7:])
	if token == "" {
		return "", errors.New("authorization header must contain a token")
	}
	return token, nil
}
