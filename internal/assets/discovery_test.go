package assets

import (
	"testing"

	"github.com/networkscope/netscope/pkg/models"
)

func TestDiscoverTarget(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		wantLen int
		wantType models.AssetType
	}{
		{"IPv4", "192.168.1.1", 1, models.AssetTypeIP},
		{"IPv6", "::1", 1, models.AssetTypeIP},
		{"domain", "example.com", 1, models.AssetTypeDomain},
		{"hostname", "web01", 1, models.AssetTypeHostname},
		{"empty", "", 0, ""},
		{"spaces", "  example.com  ", 1, models.AssetTypeDomain},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assets, err := DiscoverTarget(tt.target)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(assets) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(assets), tt.wantLen)
			}
			if tt.wantLen > 0 && assets[0].Type != tt.wantType {
				t.Errorf("type = %q, want %q", assets[0].Type, tt.wantType)
			}
		})
	}
}

func TestRegistryAddAndGet(t *testing.T) {
	r := NewRegistry()
	a, err := models.NewAsset("a1", models.AssetTypeIP, "dns")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if !r.Add(a) {
		t.Error("first Add should return true")
	}
	if r.Add(a) {
		t.Error("duplicate Add should return false")
	}
	if r.Count() != 1 {
		t.Errorf("Count = %d, want 1", r.Count())
	}
	got, ok := r.Get("a1")
	if !ok || got.ID != "a1" {
		t.Errorf("Get returned unexpected result")
	}
	_, ok = r.Get("missing")
	if ok {
		t.Error("Get should return false for missing asset")
	}
}

func TestRegistryDedup(t *testing.T) {
	r := NewRegistry()
	a1, _ := models.NewAsset("a1", models.AssetTypeIP, "dns")
	a2, _ := models.NewAsset("a1", models.AssetTypeIP, "scan")
	r.Add(a1)
	r.Add(a2)
	if r.Count() != 1 {
		t.Errorf("Count = %d, want 1 after duplicate add", r.Count())
	}
}

func FuzzDiscoverTarget(f *testing.F) {
	f.Add("192.168.1.1")
	f.Add("example.com")
	f.Add("web01")
	f.Add("::1")
	f.Fuzz(func(t *testing.T, target string) {
		_, _ = DiscoverTarget(target)
	})
}
