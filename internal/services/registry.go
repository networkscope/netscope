package services

import (
	"github.com/networkscope/netscope/pkg/models"
)

type Registry struct {
	services map[string]*models.Service
}

func NewRegistry() *Registry {
	return &Registry{services: make(map[string]*models.Service)}
}

func (r *Registry) Add(s *models.Service) bool {
	if s == nil {
		return false
	}
	if _, exists := r.services[s.ID]; exists {
		s.UpdateSeen()
		return false
	}
	r.services[s.ID] = s
	return true
}

func (r *Registry) Get(id string) (*models.Service, bool) {
	s, ok := r.services[id]
	return s, ok
}

func (r *Registry) All() []*models.Service {
	out := make([]*models.Service, 0, len(r.services))
	for _, s := range r.services {
		out = append(out, s)
	}
	return out
}

func (r *Registry) Count() int {
	return len(r.services)
}

func (r *Registry) ForAsset(assetID string) []*models.Service {
	var out []*models.Service
	for _, s := range r.services {
		if s.AssetID == assetID {
			out = append(out, s)
		}
	}
	return out
}
