package assets

import (
	"testing"

	"github.com/networkscope/netscope/pkg/models"
)

func BenchmarkDiscoverTarget(b *testing.B) {
	targets := []string{"192.168.1.1", "example.com", "web01", "::1"}
	for i := 0; i < b.N; i++ {
		for _, t := range targets {
			_, _ = DiscoverTarget(t)
		}
	}
}

func BenchmarkRegistryAdd(b *testing.B) {
	r := NewRegistry()
	a, _ := models.NewAsset("a1", models.AssetTypeIP, "dns")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Add(a)
	}
}
