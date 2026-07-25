package api

import (
	"net/http"
	"time"

	"github.com/networkscope/netscope/internal/changes"
	"github.com/networkscope/netscope/pkg/models"
)

// Server exposes the NetScope assessment engine over HTTP.
type Server struct {
	engine coreEngine
	mux    *http.ServeMux
	srv    *http.Server
}

// coreEngine is the runtime surface the API needs from the core package.
type coreEngine interface {
	Scan(string) ([]*models.Asset, error)
	Assets() []*models.Asset
	Services() []*models.Service
	Findings() []*models.Finding
	Graph() *models.Graph
	Load(string) error
	LoadPreviousSnapshot(string) (*changes.Snapshot, error)
	Snapshot(string) (*changes.Snapshot, error)
}

// ScanResponse is the stable JSON schema for scan results.
type ScanResponse struct {
	Target   string            `json:"target"`
	Assets   []*models.Asset   `json:"assets"`
	Services []*models.Service `json:"services"`
	Findings []*models.Finding `json:"findings"`
	Graph    *models.Graph     `json:"graph"`
}

// HealthResponse is the health check schema.
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// ChangesResponse is the changes schema.
type ChangesResponse struct {
	Changes []Change `json:"changes"`
	Summary Summary  `json:"summary"`
}

// Change is a single change record.
type Change struct {
	Kind    string `json:"kind"`
	Type    string `json:"type"`
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

// Summary is a count of changes by kind.
type Summary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
}
