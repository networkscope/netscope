package graph

import (
	"fmt"
	"sort"

	"github.com/networkscope/netscope/pkg/models"
)

type Builder struct {
	g *models.Graph
}

func NewBuilder() *Builder {
	return &Builder{g: models.NewGraph()}
}

func (b *Builder) Graph() *models.Graph {
	return b.g
}

func (b *Builder) AddAsset(a *models.Asset) error {
	if a == nil {
		return nil
	}
	if err := b.g.AddNode(&models.Node{ID: a.ID, Type: "asset", Label: string(a.Type)}); err != nil {
		return err
	}
	return nil
}

func (b *Builder) AddService(s *models.Service) error {
	if s == nil {
		return nil
	}
	if err := b.g.AddNode(&models.Node{ID: s.ID, Type: "service", Label: fmt.Sprintf("%s/%s", s.Transport, s.Protocol)}); err != nil {
		return err
	}
	return b.g.AddEdge(s.AssetID, s.ID, "hosts")
}

func (b *Builder) AddFinding(f *models.Finding) error {
	if f == nil {
		return nil
	}
	if err := b.g.AddNode(&models.Node{ID: f.ID, Type: "finding", Label: f.Title}); err != nil {
		return err
	}
	return b.g.AddEdge(f.AffectedAsset, f.ID, "affects")
}

func (b *Builder) Populate(assets []*models.Asset, services []*models.Service, findings []*models.Finding) error {
	for _, a := range assets {
		if err := b.AddAsset(a); err != nil {
			return err
		}
	}
	for _, s := range services {
		if err := b.AddService(s); err != nil {
			return err
		}
	}
	for _, f := range findings {
		if err := b.AddFinding(f); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) PopulateFromLoaded(assets []*models.Asset, services []*models.Service, findings []*models.Finding, nodes []*models.Node, edges []models.Edge) {
	for _, a := range assets {
		_ = b.AddAsset(a)
	}
	for _, s := range services {
		_ = b.AddService(s)
	}
	for _, f := range findings {
		_ = b.AddFinding(f)
	}
	for _, n := range nodes {
		_ = b.g.AddNode(n)
	}
	for _, e := range edges {
		_ = b.g.AddEdge(e.Source, e.Target, e.Type)
	}
}

func TextOutput(g *models.Graph) string {
	nodes := g.Nodes()
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	edges := g.Edges()

	var out string
	for _, n := range nodes {
		out += fmt.Sprintf("%s [%s] %s\n", n.ID, n.Type, n.Label)
	}
	for _, e := range edges {
		out += fmt.Sprintf("%s -> %s (%s)\n", e.Source, e.Target, e.Type)
	}
	return out
}
