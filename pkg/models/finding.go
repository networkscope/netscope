package models

import (
	"errors"
	"fmt"
	"time"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

type FindingStatus string

const (
	StatusOpen      FindingStatus = "open"
	StatusConfirmed FindingStatus = "confirmed"
	StatusFixed     FindingStatus = "fixed"
	StatusIgnored   FindingStatus = "ignored"
)

type Finding struct {
	ID            string
	Title         string
	Severity      Severity
	AffectedAsset string
	Evidence      string
	Description   string
	Recommendation string
	References    []string
	FirstSeen     time.Time
	LastSeen      time.Time
	Status        FindingStatus
	Confidence    float64
}

func NewFinding(id, title string, severity Severity, affectedAsset string) (*Finding, error) {
	if id == "" {
		return nil, errors.New("finding ID cannot be empty")
	}
	if title == "" {
		return nil, errors.New("finding title cannot be empty")
	}
	if affectedAsset == "" {
		return nil, errors.New("finding affected asset cannot be empty")
	}
	if !isValidSeverity(severity) {
		return nil, fmt.Errorf("invalid severity: %s", severity)
	}
	now := time.Now().UTC()
	return &Finding{
		ID:            id,
		Title:         title,
		Severity:      severity,
		AffectedAsset: affectedAsset,
		Status:        StatusOpen,
		FirstSeen:     now,
		LastSeen:      now,
		References:    make([]string, 0),
	}, nil
}

func isValidSeverity(s Severity) bool {
	switch s {
	case SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo:
		return true
	}
	return false
}

func (f *Finding) UpdateSeen() {
	if f != nil {
		f.LastSeen = time.Now().UTC()
	}
}

func (f *Finding) AddReference(ref string) {
	if f != nil && ref != "" {
		f.References = append(f.References, ref)
	}
}
