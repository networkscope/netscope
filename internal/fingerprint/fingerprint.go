package fingerprint

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const bannerReadTimeout = 500 * time.Millisecond
const maxBannerBytes = 256

type BannerMatch struct {
	Software    string
	Version     string
	Confidence  float64
	Banner      string
}

// ProbeBanner connects to a TCP port and reads the initial banner.
// Returns an empty string if no banner is received within the timeout.
func ProbeBanner(host string, port int, timeout time.Duration) string {
	target := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return ""
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(bannerReadTimeout))
	buf := make([]byte, maxBannerBytes)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return ""
	}
	return strings.TrimSpace(string(buf[:n]))
}

// MatchBanner compares a banner against known service signatures.
func MatchBanner(banner string) *BannerMatch {
	banner = strings.TrimSpace(banner)
	if banner == "" {
		return nil
	}

	// HTTP Server
	if strings.HasPrefix(banner, "HTTP/") {
		server := extractHeader(banner, "Server")
		if server != "" {
			return &BannerMatch{
				Software:   server,
				Confidence: 0.85,
				Banner:     banner,
			}
		}
	}

	// SSH
	if strings.HasPrefix(banner, "SSH-") {
		parts := strings.SplitN(banner, " ", 2)
		software := parts[0] // e.g. "SSH-2.0-OpenSSH_8.9"
		if idx := strings.LastIndex(software, "-"); idx != -1 {
			software = software[idx+1:]
		}
		version := ""
		if len(parts) > 1 {
			version = parts[1]
		}
		return &BannerMatch{
			Software:   software,
			Version:    version,
			Confidence: 0.95,
			Banner:     banner,
		}
	}

	// SMTP
	if strings.Contains(banner, "ESMTP") || strings.Contains(banner, "SMTP") {
		software := extractBetween(banner, " ", " ")
		version := extractVersion(banner)
		return &BannerMatch{
			Software:   software,
			Version:    version,
			Confidence: 0.8,
			Banner:     banner,
		}
	}

	// MySQL
	if strings.Contains(banner, "mysql") || strings.Contains(banner, "MariaDB") {
		return &BannerMatch{
			Software:   "mysql",
			Confidence: 0.7,
			Banner:     banner,
		}
	}

	// FTP
	if strings.HasPrefix(banner, "220 ") || strings.Contains(banner, "FTP") {
		return &BannerMatch{
			Software:   "ftp",
			Confidence: 0.6,
			Banner:     banner,
		}
	}

	return nil
}

func extractHeader(banner, key string) string {
	lines := strings.Split(banner, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, key+":") {
			return strings.TrimSpace(line[len(key)+1:])
		}
		if strings.HasPrefix(line, key+" ") {
			return strings.TrimSpace(line[len(key)+1:])
		}
	}
	// HTTP/1.1 200 ... style
	parts := strings.SplitN(banner, " ", 3)
	if len(parts) >= 3 {
		return strings.TrimSpace(parts[2])
	}
	return ""
}

func extractBetween(s, start, end string) string {
	idx := strings.Index(s, start)
	if idx == -1 {
		return ""
	}
	idx += len(start)
	endIdx := strings.Index(s[idx:], end)
	if endIdx == -1 {
		return s[idx:]
	}
	return s[idx : idx+endIdx]
}

func extractVersion(s string) string {
	parts := strings.Split(s, " ")
	for _, p := range parts {
		if len(p) > 0 && (p[0] == 'v' || p[0] == 'V') {
			return p[1:]
		}
	}
	return ""
}
