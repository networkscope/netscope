package core

import (
	"os"
	"testing"

	"github.com/networkscope/netscope/internal/changes"
)

func TestEngineScanWithSave(t *testing.T) {
	e := NewEngine()
	_, err := e.Scan("192.168.1.1")
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	f, err := os.CreateTemp("", "netscope-*.db")
	if err != nil {
		t.Fatalf("temp file failed: %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())
	if err := e.Save(f.Name()); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	e = NewEngine()
	if err := e.Load(f.Name()); err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(e.Assets()) != 1 {
		t.Errorf("loaded assets = %d, want 1", len(e.Assets()))
	}
	if len(e.Services()) != 0 {
		t.Errorf("loaded services = %d, want 0", len(e.Services()))
	}
}

func TestEngineSnapshot(t *testing.T) {
	e := NewEngine()
	_, err := e.Scan("192.168.1.1")
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	snap, err := e.Snapshot("")
	if err != nil {
		t.Fatalf("snapshot failed: %v", err)
	}
	if len(snap.Assets) != 1 {
		t.Errorf("snapshot assets = %d, want 1", len(snap.Assets))
	}
	if snap.ID == "" {
		t.Error("snapshot id is empty")
	}
}

func TestEngineServices(t *testing.T) {
	e := NewEngine()
	_, err := e.Scan("127.0.0.1")
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(e.Services()) != 1 {
		t.Errorf("services = %d, want 1", len(e.Services()))
	}
}

func TestEngineGraphPopulated(t *testing.T) {
	e := NewEngine()
	_, err := e.Scan("127.0.0.1")
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	g := e.Graph()
	nodes := g.Nodes()
	if len(nodes) == 0 {
		t.Error("graph nodes should not be empty after scan")
	}
}

func TestEngineFindingsPopulated(t *testing.T) {
	e := NewEngine()
	_, err := e.Scan("")
	if err == nil {
		t.Fatal("expected error for empty target, got nil")
	}
	if len(e.Findings()) != 0 {
		t.Errorf("findings = %d, want 0 for failed scan", len(e.Findings()))
	}
}

func TestEngineChanges(t *testing.T) {
	f, err := os.CreateTemp("", "netscope-*.db")
	if err != nil {
		t.Fatalf("temp file failed: %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())
	if err := buildDB(f.Name(), []string{"192.168.1.1"}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if err := buildDB(f.Name(), []string{"10.0.0.1"}); err != nil {
		t.Fatalf("second scan failed: %v", err)
	}
	e := NewEngine()
	if err := e.Load(f.Name()); err != nil {
		t.Fatalf("load failed: %v", err)
	}
	prev, err := e.LoadPreviousSnapshot(f.Name())
	if err != nil {
		t.Fatalf("load previous snapshot failed: %v", err)
	}
	current, err := e.Snapshot("")
	if err != nil {
		t.Fatalf("snapshot failed: %v", err)
	}
	if prev.ID == current.ID {
		t.Fatalf("previous snapshot id should differ from current snapshot id")
	}
	result := changes.Diff(prev, current)
	if len(result.Changes) != 1 {
		t.Errorf("changes = %d, want 1", len(result.Changes))
	}
}

func buildDB(path string, targets []string) error {
	e := NewEngine()
	for _, t := range targets {
		if _, err := e.Scan(t); err != nil {
			return err
		}
	}
	return e.Save(path)
}
