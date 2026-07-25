package scanner

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/networkscope/netscope/internal/fingerprint"
	"github.com/networkscope/netscope/pkg/models"
)

type ScanResult struct {
	Port   int
	Open   bool
	Banner string
}

type Scanner struct {
	timeout  time.Duration
	workers  int
	results  []ScanResult
	mu       sync.Mutex
	wg       sync.WaitGroup
	sem      chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewScanner(timeout time.Duration, workers int) *Scanner {
	if workers <= 0 {
		workers = 100
	}
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Scanner{timeout: timeout, workers: workers, ctx: ctx, cancel: cancel, sem: make(chan struct{}, workers)}
}

func (s *Scanner) Scan(host string, ports []int) []ScanResult {
	s.results = make([]ScanResult, 0, len(ports))
	for _, p := range ports {
		s.wg.Add(1)
		s.sem <- struct{}{}
		go s.probe(host, p)
	}
	s.wg.Wait()
	return s.results
}

func (s *Scanner) probe(host string, port int) {
	defer s.wg.Done()
	defer func() { <-s.sem }()
	if s.ctx.Err() != nil {
		return
	}
	target := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", target, s.timeout)
	if err != nil {
		return
	}
	banner := fingerprint.ProbeBanner(host, port, s.timeout)
	conn.Close()
	s.mu.Lock()
	s.results = append(s.results, ScanResult{Port: port, Open: true, Banner: banner})
	s.mu.Unlock()
}

func (s *Scanner) Stop() {
	s.cancel()
}

func CommonPorts() []int {
	return []int{21, 22, 23, 25, 53, 80, 110, 139, 143, 443, 445, 993, 995, 1433, 1521, 3306, 3389, 5432, 5900, 6379, 8080, 8443, 9000, 27017}
}

func ResultsToServices(assetID string, results []ScanResult) []*models.Service {
	var out []*models.Service
	for _, r := range results {
		if !r.Open {
			continue
		}
		proto := guessProtocol(r.Port)
		s, _ := models.NewService(fmt.Sprintf("%s:%d", assetID, r.Port), assetID, r.Port, "tcp")
		s.Protocol = proto
		s.Confidence = 1.0

		if r.Banner != "" {
			m := fingerprint.MatchBanner(r.Banner)
			if m != nil && m.Confidence >= 0.7 {
				s.Software = m.Software
				s.Version = m.Version
				s.Name = m.Software
				s.Confidence = m.Confidence
			}
		}
		out = append(out, s)
	}
	return out
}

func guessProtocol(port int) string {
	switch port {
	case 21:
		return "ftp"
	case 22:
		return "ssh"
	case 23:
		return "telnet"
	case 25:
		return "smtp"
	case 53:
		return "dns"
	case 80, 8080:
		return "http"
	case 110:
		return "pop3"
	case 139, 445:
		return "smb"
	case 143, 993:
		return "imap"
	case 443, 8443:
		return "https"
	case 3306:
		return "mysql"
	case 3389:
		return "rdp"
	case 5432:
		return "postgresql"
	case 6379:
		return "redis"
	case 27017:
		return "mongodb"
	}
	return ""
}
