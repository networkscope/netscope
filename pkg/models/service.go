package models

import (
	"errors"
	"fmt"
	"time"
)

type Service struct {
	ID          string
	AssetID     string
	Port        int
	Transport   string
	Protocol    string
	Name        string
	Software    string
	Version     string
	Confidence  float64
	FirstSeen   time.Time
	LastSeen    time.Time
}

func NewService(id, assetID string, port int, transport string) (*Service, error) {
	if id == "" {
		return nil, errors.New("service ID cannot be empty")
	}
	if assetID == "" {
		return nil, errors.New("service asset ID cannot be empty")
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}
	if transport == "" {
		return nil, errors.New("transport protocol cannot be empty")
	}
	now := time.Now().UTC()
	return &Service{
		ID:         id,
		AssetID:    assetID,
		Port:       port,
		Transport:  transport,
		Confidence: 1.0,
		FirstSeen:  now,
		LastSeen:   now,
	}, nil
}

func (s *Service) UpdateSeen() {
	if s != nil {
		s.LastSeen = time.Now().UTC()
	}
}
