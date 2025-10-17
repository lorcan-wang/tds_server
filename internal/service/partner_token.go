package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"tds_server/internal/config"

	"github.com/go-resty/resty/v2"
)

const (
	partnerGrantType        = "client_credentials"
	defaultTokenValidity    = time.Hour
	partnerTokenRefreshSkew = 30 * time.Second
)

type partnerTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// PartnerTokenService retrieves and caches partner access tokens in memory. PartnerTokenService 在内存中获取并缓存合作伙伴访问令牌。
type PartnerTokenService struct {
	cfg        *config.Config
	client     *resty.Client
	mu         sync.RWMutex
	token      string
	expiresAt  time.Time
	registered bool
}

// NewPartnerTokenService creates a PartnerTokenService and loads the initial token eagerly. NewPartnerTokenService 会创建服务并主动加载初始令牌。
func NewPartnerTokenService(cfg *config.Config) (*PartnerTokenService, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if cfg.TeslaClientID == "" || cfg.TeslaClientSecret == "" {
		return nil, fmt.Errorf("tesla client credentials are required")
	}
	if cfg.TeslaPartnerDomain == "" {
		return nil, fmt.Errorf("tesla partner domain is required")
	}
	if _, err := url.ParseRequestURI(cfg.TeslaAPIURL); err != nil {
		return nil, fmt.Errorf("invalid tesla api url: %w", err)
	}

	svc := &PartnerTokenService{
		cfg:    cfg,
		client: resty.New(),
	}

	if err := svc.refresh(context.Background()); err != nil {
		return nil, err
	}
	return svc, nil
}

// GetToken returns a cached partner token, refreshing it if it is close to expiry. GetToken 会返回已缓存的合作伙伴令牌，并在即将过期时刷新。
func (s *PartnerTokenService) GetToken(ctx context.Context) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if token, ok := s.cachedToken(); ok {
		return token, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isTokenValidLocked() {
		return s.token, nil
	}

	if err := s.refreshLocked(ctx); err != nil {
		return "", err
	}
	return s.token, nil
}

func (s *PartnerTokenService) cachedToken() (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isTokenValidLocked() {
		return s.token, true
	}
	return "", false
}

func (s *PartnerTokenService) refresh(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.refreshLocked(ctx)
}

func (s *PartnerTokenService) refreshLocked(ctx context.Context) error {
	payload := map[string]string{
		"grant_type":    partnerGrantType,
		"client_id":     s.cfg.TeslaClientID,
		"client_secret": s.cfg.TeslaClientSecret,
		"audience":      s.cfg.TeslaAPIURL,
	}
	if scope := s.cfg.TeslaPartnerScope; scope != "" {
		payload["scope"] = scope
	}
	var tokenResp partnerTokenResponse

	resp, err := s.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", defaultUserAgent).
		SetBody(payload).
		SetResult(&tokenResp).
		Post(s.cfg.TeslaPartnerTokenURL)
	if err != nil {
		return fmt.Errorf("request partner token: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request partner token: status %d body %s", resp.StatusCode(), resp.String())
	}
	if tokenResp.AccessToken == "" {
		return fmt.Errorf("request partner token: empty access_token in response")
	}

	now := time.Now()
	expires := now.Add(defaultTokenValidity)
	if tokenResp.ExpiresIn > 0 {
		expires = now.Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	s.token = tokenResp.AccessToken
	s.expiresAt = expires

	if !s.registered {
		if err := s.registerPartnerAccount(ctx, s.token); err != nil {
			return err
		}
		s.registered = true
	}
	return nil
}

func (s *PartnerTokenService) isTokenValidLocked() bool {
	if s.token == "" {
		return false
	}
	return time.Now().Add(partnerTokenRefreshSkew).Before(s.expiresAt)
}

func (s *PartnerTokenService) registerPartnerAccount(ctx context.Context, token string) error {
	endpoint := s.cfg.TeslaAPIURL + "/api/1/partner_accounts"
	requestBody := map[string]string{
		"domain": s.cfg.TeslaPartnerDomain,
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("User-Agent", defaultUserAgent).
		SetBody(requestBody).
		Post(endpoint)
	if err != nil {
		return fmt.Errorf("register partner account: %w", err)
	}

	if resp.StatusCode() != http.StatusOK &&
		resp.StatusCode() != http.StatusCreated &&
		resp.StatusCode() != http.StatusConflict {
		return fmt.Errorf("register partner account: status %d body %s", resp.StatusCode(), resp.String())
	}

	return nil
}
