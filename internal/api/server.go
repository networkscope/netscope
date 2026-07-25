package api

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/networkscope/netscope/internal/changes"
	"github.com/networkscope/netscope/internal/core"
)

//go:embed web/*
var webFS embed.FS

func NewServer(addr string) *Server {
	s := &Server{
		engine: core.NewEngine(),
		mux:    http.NewServeMux(),
	}
	s.mux.HandleFunc("/api/v1/health", s.healthHandler)
	s.mux.HandleFunc("/api/v1/scan", s.scanHandler)
	s.mux.HandleFunc("/api/v1/changes", s.changesHandler)
	s.mux.HandleFunc("/api/v1/assets", s.assetsHandler)
	s.mux.HandleFunc("/api/v1/services", s.servicesHandler)
	s.mux.HandleFunc("/api/v1/findings", s.findingsHandler)
	s.mux.HandleFunc("/api/v1/graph", s.graphHandler)
	s.mux.HandleFunc("/", s.serveWeb)
	s.srv = &http.Server{
		Addr:              addr,
		Handler:           corsMiddleware(s.mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	return s
}

func (s *Server) serveWeb(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	if strings.HasPrefix(path, "/web/") {
		path = strings.TrimPrefix(path, "/web/")
	}
	fullPath := filepath.Join("web", path)
	data, err := webFS.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	contentType := getContentType(path)
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

func getContentType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".ts":
		return "application/typescript; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".woff", ".woff2":
		return "font/woff2"
	default:
		return "text/plain; charset=utf-8"
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) scanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "missing target parameter", http.StatusBadRequest)
		return
	}
	_, err := s.engine.Scan(target)
	if err != nil {
		http.Error(w, fmt.Sprintf("scan failed: %v", err), http.StatusInternalServerError)
		return
	}
	resp := ScanResponse{
		Target:   target,
		Assets:   s.engine.Assets(),
		Services: s.engine.Services(),
		Findings: s.engine.Findings(),
		Graph:    s.engine.Graph(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) assetsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.engine.Assets())
}

func (s *Server) servicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.engine.Services())
}

func (s *Server) findingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.engine.Findings())
}

func (s *Server) graphHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.engine.Graph())
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok", Timestamp: time.Now().UTC()})
}

func (s *Server) changesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "missing path parameter", http.StatusBadRequest)
		return
	}
	if err := s.engine.Load(path); err != nil {
		http.Error(w, fmt.Sprintf("load failed: %v", err), http.StatusInternalServerError)
		return
	}
	prev, err := s.engine.LoadPreviousSnapshot(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("load previous snapshot failed: %v", err), http.StatusInternalServerError)
		return
	}
	current, err := s.engine.Snapshot("")
	if err != nil {
		http.Error(w, fmt.Sprintf("snapshot failed: %v", err), http.StatusInternalServerError)
		return
	}
	result := changes.Diff(prev, current)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChangesResponse{
		Changes: toAPICHanges(result.Changes),
		Summary: Summary{Added: result.Summary[changes.Added], Removed: result.Summary[changes.Removed], Modified: result.Summary[changes.Modified]},
	})
}

func toAPICHanges(in []changes.Change) []Change {
	out := make([]Change, len(in))
	for i, c := range in {
		out[i] = Change{Kind: string(c.Kind), Type: c.Type, ID: c.ID, Summary: c.Summary}
	}
	return out
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
