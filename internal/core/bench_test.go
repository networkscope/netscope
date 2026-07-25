package core

import (
	"testing"
)

func BenchmarkScan(b *testing.B) {
	e := NewEngine()
	targets := []string{"192.168.1.1", "10.0.0.1", "::1", "example.com", "web01"}
	for i := 0; i < b.N; i++ {
		for _, t := range targets {
			_, err := e.Scan(t)
			if err != nil {
				b.Fatalf("scan failed: %v", err)
			}
		}
	}
}
