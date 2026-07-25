package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/networkscope/netscope/internal/core"
)

func TestRootHelp(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("help failed: %v", err)
	}
	out := buf.String()
	for _, cmd := range []string{"scan", "assets", "services", "findings", "graph", "changes", "report"} {
		if !strings.Contains(out, cmd) {
			t.Errorf("help missing command %q", cmd)
		}
	}
}

func TestScanSubcommand(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"scan", "192.168.1.1"})
	engine := core.NewEngine()
	if _, err := engine.Scan("192.168.1.1"); err != nil {
		t.Fatalf("scan setup failed: %v", err)
	}
	if engine.Assets() == nil || len(engine.Assets()) != 1 {
		t.Fatal("engine assets not populated")
	}
}

func TestScanMissingArg(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"scan"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing scan arg, got nil")
	}
}

func TestQuietMode(t *testing.T) {
	savePath = "tmp.db"
	quiet = true
	format = "text"
	rootCmd.SetArgs([]string{"scan", "192.168.1.1"})
	engine := core.NewEngine()
	_, err := engine.Scan("192.168.1.1")
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	printResults(engine, "192.168.1.1")
	quiet = false
	savePath = ""
}
