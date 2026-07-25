package graph

import (
	"testing"

	"github.com/networkscope/netscope/pkg/models"
)

func TestBuilderPopulate(t *testing.T) {
	b := NewBuilder()
	assets, _ := models.NewAsset("a1", models.AssetTypeIP, "scan")
	svc, _ := models.NewService("s1", "a1", 22, "tcp")
	finding, _ := models.NewFinding("f1", "finding", models.SeverityLow, "a1")
	if err := b.Populate([]*models.Asset{assets}, []*models.Service{svc}, []*models.Finding{finding}); err != nil {
		t.Fatalf("Populate failed: %v", err)
	}
	nodes := b.Graph().Nodes()
	if len(nodes) != 3 {
		t.Errorf("Nodes = %d, want 3", len(nodes))
	}
	edges := b.Graph().Edges()
	if len(edges) != 2 {
		t.Errorf("Edges = %d, want 2", len(edges))
	}
	if edges[0].Type != "hosts" {
		t.Errorf("edge type = %q, want hosts", edges[0].Type)
	}
	if edges[1].Type != "affects" {
		t.Errorf("edge type = %q, want affects", edges[1].Type)
	}
}

func TestTextOutput(t *testing.T) {
	b := NewBuilder()
	a, _ := models.NewAsset("a1", models.AssetTypeIP, "scan")
	s, _ := models.NewService("s1", "a1", 80, "tcp")
	b.Populate([]*models.Asset{a}, []*models.Service{s}, nil)
	out := TextOutput(b.Graph())
	if out == "" {
		t.Error("TextOutput returned empty string")
	}
}
