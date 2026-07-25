package models

import (
	"errors"
	"time"
)

type AssetType string

const (
	AssetTypeIP       AssetType = "ip"
	AssetTypeDomain   AssetType = "domain"
	AssetTypeHostname AssetType = "hostname"
	AssetTypeApp      AssetType = "application"
)

type Asset struct {
	ID        string
	Type      AssetType
	Source    string
	Evidence  string
	FirstSeen time.Time
	LastSeen  time.Time
	Metadata  map[string]string
}

func NewAsset(id string, assetType AssetType, source string) (*Asset, error) {
	if id == "" {
		return nil, errors.New("asset ID cannot be empty")
	}
	if source == "" {
		return nil, errors.New("asset source cannot be empty")
	}
	now := time.Now().UTC()
	return &Asset{
		ID:        id,
		Type:      assetType,
		Source:    source,
		FirstSeen: now,
		LastSeen:  now,
		Metadata:  make(map[string]string),
	}, nil
}

func (a *Asset) UpdateSeen() {
	if a != nil {
		a.LastSeen = time.Now().UTC()
	}
}

func (a *Asset) SetMetadata(key, value string) {
	if a != nil && a.Metadata != nil {
		a.Metadata[key] = value
	}
}
