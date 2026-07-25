package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/networkscope/netscope/internal/core"
)

func TestHealthEndpoint(t *testing.T) {
	s := NewServer(":0")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rr := httptest.NewRecorder()
	s.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("health status = %d, want 200", rr.Code)
	}
	var body HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode health response: %v", err)
	}
	if body.Status != "ok" {
		t.Errorf("health status = %q, want ok", body.Status)
	}
}

func TestScanEndpoint(t *testing.T) {
	s := NewServer(":0")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scan?target=192.168.1.1", nil)
	rr := httptest.NewRecorder()
	s.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("scan status = %d, want 200", rr.Code)
	}
	var body ScanResponse
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode scan response: %v", err)
	}
	if body.Target != "192.168.1.1" {
		t.Errorf("target = %q, want 192.168.1.1", body.Target)
	}
	if len(body.Assets) != 1 {
		t.Errorf("assets = %d, want 1", len(body.Assets))
	}
}

func TestScanMissingTarget(t *testing.T) {
	s := NewServer(":0")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scan", nil)
	rr := httptest.NewRecorder()
	s.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("scan status = %d, want 400", rr.Code)
	}
}

func TestAssetsEndpoint(t *testing.T) {
	e := core.NewEngine()
	_, err := e.Scan("192.168.1.1")
	if err != nil {
		t.Fatalf("scan setup failed: %v", err)
	}
	s := NewServer(":0")
	s.engine = e
	req := httptest.NewRequest(http.MethodGet, "/api/v1/assets", nil)
	rr := httptest.NewRecorder()
	s.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("assets status = %d, want 200", rr.Code)
	}
	var out []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatalf("decode assets: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("assets = %d, want 1", len(out))
	}
}
