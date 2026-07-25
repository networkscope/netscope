package scanner

import (
	"testing"
	"time"
)

func TestScanLocalhost(t *testing.T) {
	s := NewScanner(200*time.Millisecond, 10)
	results := s.Scan("127.0.0.1", []int{1, 2, 65535})
	found := false
	for _, r := range results {
		if r.Open {
			found = true
			break
		}
	}
	if found {
		t.Log("found open port on localhost (expected in test environment)")
	}
}

func TestScanTimeout(t *testing.T) {
	s := NewScanner(1*time.Millisecond, 5)
	results := s.Scan("127.0.0.1", []int{1, 2, 3, 4, 5})
	for _, r := range results {
		if r.Open {
			t.Errorf("unexpected open port %d with 1ms timeout", r.Port)
		}
	}
}

func TestScanEmptyPorts(t *testing.T) {
	s := NewScanner(100*time.Millisecond, 5)
	results := s.Scan("127.0.0.1", []int{})
	if len(results) != 0 {
		t.Errorf("results = %d, want 0", len(results))
	}
}

func TestScanUnreachable(t *testing.T) {
	s := NewScanner(100*time.Millisecond, 5)
	results := s.Scan("192.0.2.1", []int{80, 443})
	for _, r := range results {
		if r.Open {
			t.Errorf("unexpected open port %d on unreachable host", r.Port)
		}
	}
}

func TestScanContextCancel(t *testing.T) {
	s := NewScanner(5*time.Second, 1)
	s.Stop()
	results := s.Scan("127.0.0.1", []int{1, 2, 3})
	_ = results
}

func TestCommonPortsCount(t *testing.T) {
	ports := CommonPorts()
	if len(ports) != 24 {
		t.Errorf("common ports = %d, want 24", len(ports))
	}
}

func TestResultsToServices(t *testing.T) {
	results := []ScanResult{
		{Port: 22, Open: true, Banner: "SSH-2.0-OpenSSH_8.9"},
		{Port: 80, Open: true, Banner: "HTTP/1.1 200 OK\r\nServer: nginx/1.18.0"},
		{Port: 443, Open: false},
	}
	out := ResultsToServices("a1", results)
	if len(out) != 2 {
		t.Errorf("services = %d, want 2", len(out))
	}
	if out[0].Port != 22 || out[1].Port != 80 {
		t.Errorf("unexpected ports: %d, %d", out[0].Port, out[1].Port)
	}
	if out[0].Software != "OpenSSH_8.9" {
		t.Errorf("software = %q, want OpenSSH_8.9", out[0].Software)
	}
	if out[1].Software != "nginx/1.18.0" {
		t.Errorf("software = %q, want nginx/1.18.0", out[1].Software)
	}
}
