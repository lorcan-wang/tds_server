package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"tds_server/internal/config"

	"github.com/go-resty/resty/v2"
)

type TeslaTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func BuildAuthURL(cfg *config.Config) string {
	params := url.Values{}
	params.Add("client_id", cfg.TeslaClientID)
	params.Add("redirect_uri", cfg.TeslaRedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "openid offline_access user_data vehicle_device_data vehicle_cmds vehicle_charging_cmds")
	params.Add("state", "db4af3f87")

	return fmt.Sprintf("%s?%s", cfg.TeslaAuthURL, params.Encode())

}

func ExchangeCode(cfg *config.Config, code string) (*TeslaTokenResponse, error) {
	client := resty.New()

	client.SetContentLength(true)
	client.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	resp, err := client.R().SetBody(map[string]interface{}{
		"grant_type":    "authorization_code",
		"client_id":     cfg.TeslaClientID,
		"client_secret": cfg.TeslaClientSecret,
		"audience":      cfg.TeslaAPIURL,
		"code":          code,
		"redirect_uri":  cfg.TeslaRedirectURI,
	}).Post(cfg.TeslaTokenURL)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange code: %d body: %s", resp.StatusCode(), string(resp.Body()))
	}

	var tr TeslaTokenResponse
	if err := json.Unmarshal(resp.Body(), &tr); err != nil {
		return nil, err
	}
	return &tr, nil
}
