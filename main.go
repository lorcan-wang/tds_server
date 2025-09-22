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
		log.Println("Warning: Error loading .env file")
	}
	teslaClientID = os.Getenv("TESLA_CLIENT_ID")
	teslaClientSecret = os.Getenv("TESLA_CLIENT_SECRET")
	teslaRedirectURI = os.Getenv("TESLA_REDIRECT_URI")
	teslaAuthURL = os.Getenv("TESLA_AUTH_URL")
	teslaTokenURL = os.Getenv("TESLA_TOKEN_URL")
	teslaAPIURL = os.Getenv("TESLA_API_URL")

	r := gin.Default()

	r.GET("/login", func(c *gin.Context) {
		oauthUrl := fmt.Sprintf("%s?&client_id=%s&redirect_uri=%s&response_type=code&scope=openid offline_access user_data vehicle_device_data vehicle_cmds vehicle_charging_cmds&state=db4af3f87", teslaAuthURL, teslaClientID, teslaRedirectURI)
		print(oauthUrl)
		c.Redirect(http.StatusTemporaryRedirect, oauthUrl)
	})

	r.GET("/auth/callback", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
			return
		}
		log.info("\n" + c)
		log.info("\n" + teslaClientID)
		log.info("\n" + teslaClientSecret)
		log.info("\n" + teslaAPIURL)
		log.info("\n" + teslaRedirectURI)
		client := resty.New()
		resp, err := client.R().SetContentLength(true).SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36").SetHeader("Host", "auth.tesla.cn").SetHeader("Content-Type", "application/json").SetBody(map[string]interface{}{
			"grant_type":    "authorization_code",
			"client_id":     teslaClientID,
			"client_secret": teslaClientSecret,
			"audience":      teslaAPIURL,
			"code":          code,
			"redirect_uri":  teslaRedirectURI,
		}).Post(teslaTokenURL)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		log.info("\n" + fmt.Sprintf("%d", resp.StatusCode()) + "\n")
		if resp.StatusCode() != http.StatusOK {
			c.JSON(resp.StatusCode(), gin.H{"error": "Failed to exchange token"})
			return
		}

		c.Data(http.StatusOK, "application/json", resp.Body())

	})
	// r.GET("/auth/callback", func(c *gin.Context) {

	// })
	r.Run(":8080")
}
