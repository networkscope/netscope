package findings

import (
	"fmt"
	"strings"

	"github.com/networkscope/netscope/internal/assets"
	"github.com/networkscope/netscope/internal/services"
	"github.com/networkscope/netscope/pkg/models"
)

type Rule struct {
	ID          string
	Description string
	Eval        func(*assets.Registry, *services.Registry) []*models.Finding
}

func defaultRules() []Rule {
	return []Rule{
		{
			ID:          "R001",
			Description: "Asset lacks discovery source detail",
			Eval: func(ar *assets.Registry, _ *services.Registry) []*models.Finding {
				var out []*models.Finding
				for _, a := range ar.All() {
					if strings.TrimSpace(a.Source) == "" {
						out = append(out, mustFinding("R001-"+a.ID, 
							"Asset missing discovery source", models.SeverityLow, a.ID,
							"source field is empty", "Record the discovery source for this asset", nil))
					}
				}
				return out
			},
		},
		{
			ID:          "R002",
			Description: "Service has low confidence identification",
			Eval: func(_ *assets.Registry, sr *services.Registry) []*models.Finding {
				var out []*models.Finding
				for _, s := range sr.All() {
					if s.Confidence < 0.5 {
						out = append(out, mustFinding("R002-"+s.ID,
							"Low confidence service identification", models.SeverityInfo, s.AssetID,
							"service confidence below 0.5", "Re-probe with additional checks", nil))
					}
				}
				return out
			},
		},
		{
			ID:          "R003",
			Description: "High-risk port exposed",
			Eval: func(_ *assets.Registry, sr *services.Registry) []*models.Finding {
				var out []*models.Finding
				highRisk := map[int]string{
					23: "telnet",
					3389: "rdp",
				}
				for _, s := range sr.All() {
					if label, ok := highRisk[s.Port]; ok {
						out = append(out, mustFinding("R003-"+s.ID,
							"High-risk port exposed: "+label, models.SeverityHigh, s.AssetID,
							fmt.Sprintf("port %d/%s open", s.Port, label), "Restrict access or disable the service", []string{"https://owasp.org/www-project-top-ten/"}))
					}
				}
				return out
			},
		},
		{
			ID:          "R004",
			Description: "HTTP service without HTTPS counterpart",
			Eval: func(_ *assets.Registry, sr *services.Registry) []*models.Finding {
				var out []*models.Finding
				for _, s := range sr.All() {
					if s.Protocol == "http" {
						hasHTTPS := false
						for _, other := range sr.All() {
							if other.Protocol == "https" && other.AssetID == s.AssetID {
								hasHTTPS = true
								break
							}
						}
						if !hasHTTPS {
							out = append(out, mustFinding("R004-"+s.ID,
								"HTTP service without HTTPS", models.SeverityMedium, s.AssetID,
								fmt.Sprintf("port %d exposes HTTP without HTTPS", s.Port), "Enable HTTPS on this asset", nil))
						}
					}
				}
				return out
			},
		},
		{
			ID:          "R005",
			Description: "Database service exposed",
			Eval: func(_ *assets.Registry, sr *services.Registry) []*models.Finding {
				var out []*models.Finding
				dbProtos := map[string]bool{"mysql": true, "postgresql": true, "mongodb": true, "redis": true}
				for _, s := range sr.All() {
					if dbProtos[s.Protocol] {
						out = append(out, mustFinding("R005-"+s.ID,
							"Database service exposed", models.SeverityHigh, s.AssetID,
							fmt.Sprintf("port %d/%s database exposed", s.Port, s.Protocol), "Restrict database access to authorized networks", []string{"https://owasp.org/www-project-top-ten/"}))
					}
				}
				return out
			},
		},
	}
}

type Evaluator struct {
	rules []Rule
	reg   *Registry
}

func NewEvaluator(reg *Registry) *Evaluator {
	return &Evaluator{
		rules: defaultRules(),
		reg:   reg,
	}
}

func (e *Evaluator) Evaluate(ar *assets.Registry, sr *services.Registry) {
	for _, rule := range e.rules {
		for _, f := range rule.Eval(ar, sr) {
			e.reg.Add(f)
		}
	}
}

func mustFinding(id, title string, sev models.Severity, asset, evidence, rec string, refs []string) *models.Finding {
	f, err := models.NewFinding(id, title, sev, asset)
	if err != nil {
		return nil
	}
	f.Evidence = evidence
	f.Recommendation = rec
	if refs != nil {
		for _, ref := range refs {
			f.AddReference(ref)
		}
	}
	f.Confidence = 1.0
	return f
}
