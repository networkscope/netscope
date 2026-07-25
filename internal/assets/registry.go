package assets

import (
	"github.com/networkscope/netscope/pkg/models"
)

type Registry struct {
	assets map[string]*models.Asset
}

func NewRegistry() *Registry {
	return &Registry{assets: make(map[string]*models.Asset)}
}

func (r *Registry) Add(a *models.Asset) bool {
	if a == nil {
		return false
	}
	if _, exists := r.assets[a.ID]; exists {
		a.UpdateSeen()
		return false
	}
	r.assets[a.ID] = a
	return true
}

func (r *Registry) Get(id string) (*models.Asset, bool) {
	a, ok := r.assets[id]
	return a, ok
}

func (r *Registry) All() []*models.Asset {
	out := make([]*models.Asset, 0, len(r.assets))
	for _, a := range r.assets {
		out = append(out, a)
	}
	return out
}

func (r *Registry) Count() int {
	return len(r.assets)
}
