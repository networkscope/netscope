package models

import (
	"testing"
	"time"
)

func TestNewAsset(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		assetType AssetType
		source  string
		wantErr bool
	}{
		{"valid asset", "asset-1", AssetTypeIP, "dns", false},
		{"empty id", "", AssetTypeIP, "dns", true},
		{"empty source", "asset-1", AssetTypeIP, "", true},
		{"domain type", "dp1", AssetTypeDomain, "cert", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAsset(tt.id, tt.assetType, tt.source)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if a.ID != tt.id {
				t.Errorf("ID = %q, want %q", a.ID, tt.id)
			}
			if a.Type != tt.assetType {
				t.Errorf("Type = %q, want %q", a.Type, tt.assetType)
			}
			if a.Source != tt.source {
				t.Errorf("Source = %q, want %q", a.Source, tt.source)
			}
			if a.Metadata == nil {
				t.Errorf("Metadata is nil")
			}
		})
	}
}

func TestAssetUpdateSeen(t *testing.T) {
	a, err := NewAsset("a1", AssetTypeIP, "dns")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	first := a.LastSeen
	time.Sleep(2 * time.Millisecond)
	a.UpdateSeen()
	if !a.LastSeen.After(first) {
		t.Errorf("LastSeen was not updated")
	}
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		assetID   string
		port      int
		transport string
		wantErr   bool
		errMsg    string
	}{
		{"valid service", "s1", "a1", 443, "tcp", false, ""},
		{"empty id", "", "a1", 443, "tcp", true, "service ID cannot be empty"},
		{"empty asset", "s1", "", 443, "tcp", true, "service asset ID cannot be empty"},
		{"invalid port low", "s1", "a1", 0, "tcp", true, "port must be between 1 and 65535"},
		{"invalid port high", "s1", "a1", 65536, "tcp", true, "port must be between 1 and 65535"},
		{"empty transport", "s1", "a1", 22, "", true, "transport protocol cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewService(tt.id, tt.assetID, tt.port, tt.transport)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errMsg)
				}
				if tt.errMsg != "" && !stringsContain(err.Error(), tt.errMsg) {
					t.Errorf("error = %q, want containing %q", err.Error(), tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s.Port != tt.port {
				t.Errorf("Port = %d, want %d", s.Port, tt.port)
			}
			if s.Confidence != 1.0 {
				t.Errorf("Confidence = %f, want 1.0", s.Confidence)
			}
		})
	}
}

func stringsContain(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestNewFinding(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		title         string
		severity      Severity
		affectedAsset string
		wantErr       bool
	}{
		{"valid finding", "f1", "Open SSH", SeverityMedium, "a1", false},
		{"empty id", "", "title", SeverityMedium, "a1", true},
		{"empty title", "f1", "", SeverityMedium, "a1", true},
		{"empty affected", "f1", "title", SeverityMedium, "", true},
		{"invalid severity", "f1", "title", "bad", "a1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFinding(tt.id, tt.title, tt.severity, tt.affectedAsset)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f.Severity != tt.severity {
				t.Errorf("Severity = %q, want %q", f.Severity, tt.severity)
			}
			if f.Status != StatusOpen {
				t.Errorf("Status = %q, want %q", f.Status, StatusOpen)
			}
		})
	}
}

func TestFindingAddReference(t *testing.T) {
	f, err := NewFinding("f1", "title", SeverityLow, "a1")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	f.AddReference("https://example.com")
	if len(f.References) != 1 || f.References[0] != "https://example.com" {
		t.Errorf("References = %v, want [https://example.com]", f.References)
	}
}

func TestGraph(t *testing.T) {
	g := NewGraph()

	node1 := &Node{ID: "a1", Type: "asset", Label: "Asset 1"}
	node2 := &Node{ID: "s1", Type: "service", Label: "SSH"}

	if err := g.AddNode(node1); err != nil {
		t.Fatalf("AddNode failed: %v", err)
	}
	if err := g.AddNode(node2); err != nil {
		t.Fatalf("AddNode failed: %v", err)
	}
	if err := g.AddNode(node1); err != nil {
		t.Fatalf("re-add should not error: %v", err)
	}

	err := g.AddEdge("a1", "s1", "hosts")
	if err != nil {
		t.Fatalf("AddEdge failed: %v", err)
	}

	if len(g.Nodes()) != 2 {
		t.Errorf("Nodes count = %d, want 2", len(g.Nodes()))
	}
	if len(g.Edges()) != 1 {
		t.Errorf("Edges count = %d, want 1", len(g.Edges()))
	}
	if g.Edges()[0].Type != "hosts" {
		t.Errorf("Edge type = %q, want %q", g.Edges()[0].Type, "hosts")
	}

	badNode := &Node{ID: "", Type: "x"}
	if err := g.AddNode(badNode); err == nil {
		t.Error("expected error for empty node ID, got nil")
	}

	if err := g.AddEdge("", "s1", "hosts"); err == nil {
		t.Error("expected error for empty source, got nil")
	}
	if err := g.AddEdge("a1", "", "hosts"); err == nil {
		t.Error("expected error for empty target, got nil")
	}
	if err := g.AddEdge("a1", "s1", ""); err == nil {
		t.Error("expected error for empty edge type, got nil")
	}
}
