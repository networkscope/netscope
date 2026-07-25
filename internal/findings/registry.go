package findings

import (
	"github.com/networkscope/netscope/pkg/models"
)

type Registry struct {
	findings map[string]*models.Finding
}

func NewRegistry() *Registry {
	return &Registry{findings: make(map[string]*models.Finding)}
}

func (r *Registry) Add(f *models.Finding) bool {
	if f == nil {
		return false
	}
	if _, exists := r.findings[f.ID]; exists {
		f.UpdateSeen()
		return false
	}
	r.findings[f.ID] = f
	return true
}

func (r *Registry) Get(id string) (*models.Finding, bool) {
	f, ok := r.findings[id]
	return f, ok
}

func (r *Registry) All() []*models.Finding {
	out := make([]*models.Finding, 0, len(r.findings))
	for _, f := range r.findings {
		out = append(out, f)
	}
	return out
}

func (r *Registry) Count() int {
	return len(r.findings)
}
