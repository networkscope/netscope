package services

import (
	"fmt"
	"time"

	"github.com/networkscope/netscope/internal/scanner"
	"github.com/networkscope/netscope/pkg/models"
)

// AnalyzeTarget probes the target for common services using TCP connect scanning.
// It returns services for open ports with basic protocol inference.
func AnalyzeTarget(assetID string) ([]*models.Service, error) {
	if assetID == "" {
		return nil, fmt.Errorf("asset ID cannot be empty")
	}
	ports := scanner.CommonPorts()
	s := scanner.NewScanner(2*time.Second, 100)
	results := s.Scan(assetID, ports)
	return scanner.ResultsToServices(assetID, results), nil
}
