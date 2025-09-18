package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

var (
	teslaClientID     string
	teslaClientSecret string
	teslaRedirectURI  string
	teslaAuthURL      string
	teslaTokenURL     string
	teslaAPIURL       string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	teslaClientID = os.Getenv("TESLA_CLIENT_ID")
	teslaClientSecret = os.Getenv("TESLA_CLIENT_SECRET")
	teslaRedirectURI = os.Getenv("TESLA_REDIRECT_URI")
	teslaAuthURL = os.Getenv("TESLA_AUTH_URL")
	teslaTokenURL = os.Getenv("TESLA_TOKEN_URL")
	teslaAPIURL = os.Getenv("TESLA_API_URL")

	r := gin.Default()

	r.GET("/login", func(c *gin.Context) {
		oauthUrl := fmt.Sprintf("%s?&client_id=%s&redirect_uri=%s&response_type=code&scope=openid offline_access user_data vehicle_device_data vehicle_cmds vehicle_charging_cmds&state=123123", teslaAuthURL, teslaClientID, teslaRedirectURI)
		c.Redirect(http.StatusTemporaryRedirect, oauthUrl)
	})

	r.GET("/auth/callback", func(c *gin.Context) {
		print("11132gag")
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
			return
		}

		client := resty.New()
		resp, err := client.R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     teslaClientID,
			"client_secret": teslaClientSecret,
			"audience":      teslaAPIURL,
			"code":          code,
			"state":         "123123",
			"redirect_uri":  "http://localhost:3000",
			"scope":         "openid offline_access user_data vehicle_device_data vehicle_cmds vehicle_charging_cmds",
		}).Post(teslaTokenURL)

		if err != nil {
			print("1gagaga")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		print(string(resp.Body()))
		c.Data(http.StatusOK, "application/json", resp.Body())

	})
	// r.GET("/auth/callback", func(c *gin.Context) {

	// })
	r.Run(":3000")
}
