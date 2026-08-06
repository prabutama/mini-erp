package httpfiber

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/service"
)

func TestSignup(t *testing.T) {
	app := NewRouter(service.NewAuthService())
	body, _ := json.Marshal(domain.SignupRequest{
		BusinessName:     "Acme Services",
		BusinessTimezone: "Asia/Jakarta",
		AdminFullName:    "Owner",
		AdminEmail:       "owner@acme.test",
		AdminPassword:    "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestLoginAndMe(t *testing.T) {
	auth := service.NewAuthService()
	app := NewRouter(auth)

	loginBody, _ := json.Marshal(domain.LoginRequest{Email: "platform.admin@example.com", Password: "admin123"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, err := app.Test(loginReq)
	if err != nil {
		t.Fatal(err)
	}
	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", loginResp.StatusCode)
	}

	var session domain.AuthSession
	if err := json.NewDecoder(loginResp.Body).Decode(&session); err != nil {
		t.Fatal(err)
	}

	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+session.AccessToken)
	meResp, err := app.Test(meReq)
	if err != nil {
		t.Fatal(err)
	}
	if meResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", meResp.StatusCode)
	}
}

func TestMeUnauthorized(t *testing.T) {
	app := NewRouter(service.NewAuthService())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}
