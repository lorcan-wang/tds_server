package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"tds_server/internal/config"
	"tds_server/internal/repository"

	"github.com/gin-gonic/gin"
)

func TestPublicKeyEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := NewRouter(&config.Config{}, &repository.TokenRepo{})

	req := httptest.NewRequest(http.MethodGet, "/.well-known/appspecific/com.tesla.3p.public-key.pem", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d want %d", rec.Code, http.StatusOK)
	}

	expectedPath := publicKeyFilePath()
	expected, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read expected public key: %v", err)
	}

	if !bytes.Equal(rec.Body.Bytes(), expected) {
		t.Fatalf("served public key does not match file contents")
	}
}
