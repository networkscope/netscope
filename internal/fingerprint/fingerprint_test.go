package fingerprint

import (
	"testing"
	"time"
)

func TestProbeBanner(t *testing.T) {
	banner := ProbeBanner("127.0.0.1", 1, 100*time.Millisecond)
	_ = banner
}

func TestMatchSSHBanner(t *testing.T) {
	m := MatchBanner("SSH-2.0-OpenSSH_8.9p1 Ubuntu-3ubuntu0.3")
	if m == nil {
		t.Fatal("expected match for SSH banner")
	}
	if m.Software != "OpenSSH_8.9p1" {
		t.Errorf("software = %q, want OpenSSH_8.9p1", m.Software)
	}
	if m.Confidence < 0.9 {
		t.Errorf("confidence = %f, want >= 0.9", m.Confidence)
	}
}

func TestMatchHTTPBanner(t *testing.T) {
	m := MatchBanner("HTTP/1.1 200 OK\r\nServer: nginx/1.18.0\r\nContent-Type: text/html")
	if m == nil {
		t.Fatal("expected match for HTTP banner")
	}
	if m.Software != "nginx/1.18.0" {
		t.Errorf("software = %q, want nginx/1.18.0", m.Software)
	}
	if m.Confidence < 0.8 {
		t.Errorf("confidence = %f, want >= 0.8", m.Confidence)
	}
}

func TestMatchSMTPBanner(t *testing.T) {
	m := MatchBanner("220 mail.example.com ESMTP Postfix (Ubuntu)")
	if m == nil {
		t.Fatal("expected match for SMTP banner")
	}
	if m.Confidence < 0.7 {
		t.Errorf("confidence = %f, want >= 0.7", m.Confidence)
	}
}

func TestMatchNoBanner(t *testing.T) {
	m := MatchBanner("")
	if m != nil {
		t.Error("expected no match for empty banner")
	}
}

func TestMatchUnknownBanner(t *testing.T) {
	m := MatchBanner("random data with no pattern")
	if m != nil {
		t.Errorf("expected no match for unknown banner, got %v", m)
	}
}
