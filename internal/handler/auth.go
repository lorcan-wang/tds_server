package handler

import (
	"errors"
	"net/http"
	"tds_server/internal/config"
	"tds_server/internal/repository"
	"tds_server/internal/service"
	"time"

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

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID.String(),
			"jwt": gin.H{
				"token":      jwtToken,
				"expires_in": int(cfg.JWT.Expiration.Seconds()),
				"issuer":     cfg.JWT.Issuer,
			},
			"tesla_token": teslaTokenRepo,
		})
	}
}

func resolveUserID(state string) (uuid.UUID, error) {
	if state == "" {
		return uuid.New(), nil
	}

	id, err := uuid.Parse(state)
	if err != nil {
		return uuid.Nil, errors.New("state must be a valid UUID")
	}
	return id, nil
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
