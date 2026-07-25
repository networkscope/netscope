package cli

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/networkscope/netscope/internal/core"
)

func TestPrintResultsJSONStable(t *testing.T) {
	engine := core.NewEngine()
	_, err := engine.Scan("192.168.1.1")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	format = "json"
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	printResults(engine, "192.168.1.1")
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()
	var body map[string]interface{}
	if err := json.Unmarshal([]byte(out), &body); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	if body["target"] != "192.168.1.1" {
		t.Errorf("target = %v, want 192.168.1.1", body["target"])
	}
	assets, ok := body["assets"].([]interface{})
	if !ok {
		t.Fatalf("assets type = %T, want []", body["assets"])
	}
	if len(assets) != 1 {
		t.Errorf("assets = %d, want 1", len(assets))
	}
}

func TestPrintResultsCSV(t *testing.T) {
	engine := core.NewEngine()
	_, err := engine.Scan("192.168.1.1")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	format = "csv"
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	printResults(engine, "192.168.1.1")
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()
	reader := csv.NewReader(strings.NewReader(out))
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("invalid csv: %v\n%s", err, out)
	}
	if len(rows) < 2 {
		t.Errorf("csv rows = %d, want at least 2", len(rows))
	}
	header := rows[0]
	if header[0] != "type" || header[1] != "id" {
		t.Errorf("csv header = %v, want [type id parent extra]", header)
	}
}
