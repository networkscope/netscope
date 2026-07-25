package findings

import (
	"testing"

	"github.com/networkscope/netscope/internal/assets"
	"github.com/networkscope/netscope/internal/services"
	"github.com/networkscope/netscope/pkg/models"
)

func TestEvaluatorR001(t *testing.T) {
	ar := assets.NewRegistry()
	a, _ := models.NewAsset("a1", models.AssetTypeIP, "input")
	a.Source = ""
	ar.Add(a)
	reg := NewRegistry()
	eval := NewEvaluator(reg)
	eval.Evaluate(ar, services.NewRegistry())
	if reg.Count() != 1 {
		t.Fatalf("findings = %d, want 1", reg.Count())
	}
	f, _ := reg.Get("R001-a1")
	if f.Title != "Asset missing discovery source" {
		t.Errorf("title = %q, want Asset missing discovery source", f.Title)
	}
	if f.Severity != models.SeverityLow {
		t.Errorf("severity = %q, want low", f.Severity)
	}
}

func TestEvaluatorR002(t *testing.T) {
	sr := services.NewRegistry()
	s, _ := models.NewService("s1", "a1", 80, "tcp")
	s.Confidence = 0.3
	sr.Add(s)
	reg := NewRegistry()
	eval := NewEvaluator(reg)
	eval.Evaluate(assets.NewRegistry(), sr)
	if reg.Count() != 1 {
		t.Fatalf("findings = %d, want 1", reg.Count())
	}
	f, _ := reg.Get("R002-s1")
	if f == nil {
		t.Fatal("expected R002-s1 finding")
	}
}

func TestEvaluatorNoFindings(t *testing.T) {
	reg := NewRegistry()
	eval := NewEvaluator(reg)
	eval.Evaluate(assets.NewRegistry(), services.NewRegistry())
	if reg.Count() != 0 {
		t.Errorf("findings = %d, want 0", reg.Count())
	}
}

func TestEvaluatorR003HighRiskPort(t *testing.T) {
	sr := services.NewRegistry()
	s, _ := models.NewService("s1", "a1", 3389, "tcp")
	sr.Add(s)
	reg := NewRegistry()
	eval := NewEvaluator(reg)
	eval.Evaluate(assets.NewRegistry(), sr)
	if reg.Count() != 1 {
		t.Fatalf("findings = %d, want 1", reg.Count())
	}
	f, _ := reg.Get("R003-s1")
	if f == nil {
		t.Fatal("expected R003-s1 finding")
	}
	if f.Severity != models.SeverityHigh {
		t.Errorf("severity = %q, want high", f.Severity)
	}
}

func TestEvaluatorR004HTTPWithoutHTTPS(t *testing.T) {
	sr := services.NewRegistry()
	s, _ := models.NewService("s1", "a1", 80, "tcp")
	s.Protocol = "http"
	sr.Add(s)
	reg := NewRegistry()
	eval := NewEvaluator(reg)
	eval.Evaluate(assets.NewRegistry(), sr)
	if reg.Count() != 1 {
		t.Fatalf("findings = %d, want 1", reg.Count())
	}
	f, _ := reg.Get("R004-s1")
	if f == nil {
		t.Fatal("expected R004-s1 finding")
	}
}

func TestEvaluatorR005DatabaseExposed(t *testing.T) {
	sr := services.NewRegistry()
	s, _ := models.NewService("s1", "a1", 3306, "tcp")
	s.Protocol = "mysql"
	sr.Add(s)
	reg := NewRegistry()
	eval := NewEvaluator(reg)
	eval.Evaluate(assets.NewRegistry(), sr)
	if reg.Count() != 1 {
		t.Fatalf("findings = %d, want 1", reg.Count())
	}
	f, _ := reg.Get("R005-s1")
	if f == nil {
		t.Fatal("expected R005-s1 finding")
	}
	if f.Severity != models.SeverityHigh {
		t.Errorf("severity = %q, want high", f.Severity)
	}
}
