package cli

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/networkscope/netscope/internal/changes"
	"github.com/networkscope/netscope/internal/core"
	"github.com/networkscope/netscope/pkg/models"
)

func printResults(engine *core.Engine, target string) {
	assets := engine.Assets()
	svcs := engine.Services()
	findings := engine.Findings()
	g := engine.Graph()
	switch format {
	case "json":
		printJSONResults(engine, target, assets, svcs, findings, g)
		return
	case "csv":
		printCSVResults(engine, target, assets, svcs, findings, g)
		return
	}
	fmt.Printf("NetScope Assessment\n\n")
	fmt.Printf("Target: %s\n\n", target)

	printAssets(assets)
	printServices(svcs)
	printFindings(findings)
	printGraph(g)
}

func printJSONResults(engine *core.Engine, target string, assets []*models.Asset, svcs []*models.Service, findings []*models.Finding, g *models.Graph) {
	type scanOut struct {
		Target   string            `json:"target"`
		Assets   []*models.Asset   `json:"assets"`
		Services []*models.Service `json:"services"`
		Findings []*models.Finding `json:"findings"`
		Graph    *models.Graph     `json:"graph"`
	}
	out, err := json.MarshalIndent(scanOut{Target: target, Assets: assets, Services: svcs, Findings: findings, Graph: g}, "", "  ")
	if err != nil {
		printErr("marshal json: " + err.Error())
		os.Exit(1)
	}
	fmt.Println(string(out))
}

func printCSVResults(engine *core.Engine, target string, assets []*models.Asset, svcs []*models.Service, findings []*models.Finding, g *models.Graph) {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	w.Write([]string{"type", "id", "parent", "extra"})
	for _, a := range assets {
		w.Write([]string{"asset", a.ID, "", string(a.Type)})
	}
	for _, s := range svcs {
		w.Write([]string{"service", s.ID, s.AssetID, fmt.Sprintf("%s/%s/%d", s.Transport, s.Protocol, s.Port)})
	}
	for _, f := range findings {
		w.Write([]string{"finding", f.ID, f.AffectedAsset, fmt.Sprintf("%s|%s", f.Severity, f.Title)})
	}
}

func printAssets(as []*models.Asset) {
	fmt.Printf("Assets (%d)\n", len(as))
	if len(as) == 0 {
		fmt.Printf("  None discovered\n\n")
		return
	}
	groups := make(map[string][]*models.Asset)
	for _, a := range as {
		cat := assetCategory(a.Type)
		groups[cat] = append(groups[cat], a)
	}
	order := []string{"IP Addresses", "Domains", "Hostnames", "Applications"}
	for _, cat := range order {
		list := groups[cat]
		if len(list) == 0 {
			continue
		}
		fmt.Printf("  %s\n", cat)
		for _, a := range list {
			fmt.Printf("    %s\n", a.ID)
			if a.Source != "" {
				fmt.Printf("      source: %s\n", a.Source)
			}
			if a.Evidence != "" {
				fmt.Printf("      evidence: %s\n", a.Evidence)
			}
			if len(a.Metadata) > 0 {
				keys := make([]string, 0)
				for k := range a.Metadata {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					fmt.Printf("      %s=%s\n", k, a.Metadata[k])
				}
			}
		}
	}
	fmt.Println()
}

func assetCategory(t models.AssetType) string {
	switch t {
	case models.AssetTypeIP:
		return "IP Addresses"
	case models.AssetTypeDomain:
		return "Domains"
	case models.AssetTypeHostname:
		return "Hostnames"
	case models.AssetTypeApp:
		return "Applications"
	}
	return string(t)
}

func printServices(ss []*models.Service) {
	fmt.Printf("Services (%d)\n", len(ss))
	if len(ss) == 0 {
		fmt.Printf("  None discovered\n\n")
		return
	}
	sort.Slice(ss, func(i, j int) bool {
		if ss[i].AssetID != ss[j].AssetID {
			return ss[i].AssetID < ss[j].AssetID
		}
		return ss[i].Port < ss[j].Port
	})
	for _, s := range ss {
		proto := s.Protocol
		if proto == "" {
			proto = "unknown"
		}
		detail := ""
		if s.Software != "" {
			detail = "  " + s.Software
			if s.Version != "" {
				detail += " " + s.Version
			}
		}
		fmt.Printf("  %s:%d  %s/%s  confidence=%.0f%%%s\n", s.AssetID, s.Port, s.Transport, proto, s.Confidence*100, detail)
	}
	fmt.Println()
}

func printFindings(ff []*models.Finding) {
	fmt.Printf("Findings (%d)\n", len(ff))
	if len(ff) == 0 {
		fmt.Printf("  None\n\n")
		return
	}
	bySev := map[models.Severity][]*models.Finding{}
	order := []models.Severity{models.SeverityCritical, models.SeverityHigh, models.SeverityMedium, models.SeverityLow, models.SeverityInfo}
	for _, sev := range order {
		bySev[sev] = nil
	}
	for _, f := range ff {
		bySev[f.Severity] = append(bySev[f.Severity], f)
	}
	for _, sev := range order {
		list := bySev[sev]
		if len(list) == 0 {
			continue
		}
		fmt.Printf("  %s (%d)\n", strings.ToUpper(string(sev)), len(list))
		for _, f := range list {
			fmt.Printf("    %s: %s\n", f.ID, f.Title)
			if f.Evidence != "" {
				fmt.Printf("      Evidence: %s\n", f.Evidence)
			}
			if f.Description != "" {
				fmt.Printf("      Description: %s\n", f.Description)
			}
			if f.Recommendation != "" {
				fmt.Printf("      Recommendation: %s\n", f.Recommendation)
			}
		}
	}
	fmt.Println()
}

func printGraph(g *models.Graph) {
	nodes := g.Nodes()
	if len(nodes) == 0 {
		return
	}
	fmt.Printf("Graph\n")
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	for _, n := range nodes {
		fmt.Printf("  %s [%s] %s\n", n.ID, n.Type, n.Label)
	}
	edges := g.Edges()
	for _, e := range edges {
		fmt.Printf("  %s -> %s (%s)\n", e.Source, e.Target, e.Type)
	}
}

func printChanges(r *changes.ChangeResult) {
	if format == "json" {
		out, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			printErr("marshal changes: " + err.Error())
			os.Exit(1)
		}
		fmt.Println(string(out))
		return
	}
	fmt.Printf("Changes\n")
	fmt.Printf("  Added:    %d\n", r.Summary[changes.Added])
	fmt.Printf("  Removed:  %d\n", r.Summary[changes.Removed])
	fmt.Printf("  Modified: %d\n", r.Summary[changes.Modified])
	for _, c := range r.Changes {
		symbol := "?"
		switch c.Kind {
		case changes.Added:
			symbol = "+"
		case changes.Removed:
			symbol = "-"
		case changes.Modified:
			symbol = "~"
		}
		fmt.Printf("  %s %s %s: %s\n", symbol, c.Type, c.ID, c.Summary)
	}
}
