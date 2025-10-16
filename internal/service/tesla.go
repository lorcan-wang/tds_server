package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tds_server/internal/config"

	"github.com/go-resty/resty/v2"
)

const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"

type TeslaTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
}

func BuildAuthURL(cfg *config.Config) string {
	return fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=%s&scope=%s",
		cfg.TeslaAuthURL,
		cfg.TeslaClientID,
		cfg.TeslaRedirectURI,
		"code",
		cfg.TeslaPartnerScope,
	)
}

func ExchangeCode(cfg *config.Config, code string) (*TeslaTokenResponse, error) {
	client := resty.New()

	client.SetContentLength(true)
	client.SetHeader("User-Agent", defaultUserAgent)
	resp, err := client.R().SetBody(map[string]interface{}{
		"grant_type":    "authorization_code",
		"client_id":     cfg.TeslaClientID,
		"client_secret": cfg.TeslaClientSecret,
		"audience":      cfg.TeslaAPIURL,
		"code":          code,
		"redirect_uri":  cfg.TeslaRedirectURI,
		"scope":         cfg.TeslaPartnerScope,
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
	fmt.Printf("Tesla token scopes: %s\n", tr.Scope)
	fmt.Printf("Response: %+s\n", resp.Body())
	return &tr, nil
}

func RefreshToken(cfg *config.Config, refreshToken string) (*TeslaTokenResponse, error) {
	client := resty.New()

	client.SetContentLength(true)
	client.SetHeader("User-Agent", defaultUserAgent)
	resp, err := client.R().SetBody(map[string]interface{}{
		"grant_type":    "refresh_token",
		"client_id":     cfg.TeslaClientID,
		"client_secret": cfg.TeslaClientSecret,
		"audience":      cfg.TeslaAPIURL,
		"refresh_token": refreshToken,
	}).Post(cfg.TeslaTokenURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to refresh token: %d body: %s", resp.StatusCode(), string(resp.Body()))
	}

	var tr TeslaTokenResponse
	if err := json.Unmarshal(resp.Body(), &tr); err != nil {
		return nil, err
	}
	return &tr, nil
}
