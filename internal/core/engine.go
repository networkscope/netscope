package core

import (
	"errors"
	"fmt"

	"github.com/networkscope/netscope/internal/assets"
	"github.com/networkscope/netscope/internal/changes"
	"github.com/networkscope/netscope/internal/findings"
	"github.com/networkscope/netscope/internal/graph"
	"github.com/networkscope/netscope/internal/services"
	"github.com/networkscope/netscope/internal/storage"
	"github.com/networkscope/netscope/pkg/models"
)

type Engine struct {
	registry     *assets.Registry
	svcRegistry  *services.Registry
	findReg      *findings.Registry
	builder      *graph.Builder
}

func NewEngine() *Engine {
	fr := findings.NewRegistry()
	return &Engine{
		registry:     assets.NewRegistry(),
		svcRegistry:  services.NewRegistry(),
		findReg:      fr,
		builder:      graph.NewBuilder(),
	}
}

func (e *Engine) Scan(target string) ([]*models.Asset, error) {
	if target == "" {
		return nil, errors.New("target cannot be empty")
	}
	found, err := assets.DiscoverTarget(target)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}
	for _, a := range found {
		e.registry.Add(a)
		svcs, err := services.AnalyzeTarget(a.ID)
		if err != nil {
			continue
		}
		for _, s := range svcs {
			e.svcRegistry.Add(s)
		}
	}
	evaluator := findings.NewEvaluator(e.findReg)
	evaluator.Evaluate(e.registry, e.svcRegistry)

	e.builder.Populate(e.registry.All(), e.svcRegistry.All(), e.findReg.All())
	return e.registry.All(), nil
}

func (e *Engine) Assets() []*models.Asset {
	return e.registry.All()
}

func (e *Engine) Services() []*models.Service {
	return e.svcRegistry.All()
}

func (e *Engine) Findings() []*models.Finding {
	return e.findReg.All()
}

func (e *Engine) Graph() *models.Graph {
	return e.builder.Graph()
}

func (e *Engine) Save(path string) error {
	st, err := storage.Open(path)
	if err != nil {
		return fmt.Errorf("open store failed: %w", err)
	}
	defer st.Close()
	for _, a := range e.registry.All() {
		if err := st.SaveAsset(a); err != nil {
			return fmt.Errorf("save asset failed: %w", err)
		}
	}
	for _, s := range e.svcRegistry.All() {
		if err := st.SaveService(s); err != nil {
			return fmt.Errorf("save service failed: %w", err)
		}
	}
	for _, f := range e.findReg.All() {
		if err := st.SaveFinding(f); err != nil {
			return fmt.Errorf("save finding failed: %w", err)
		}
	}
	if err := st.SaveGraphNodes(e.builder.Graph().Nodes()); err != nil {
		return fmt.Errorf("save graph nodes failed: %w", err)
	}
	if err := st.SaveGraphEdges(e.builder.Graph().Edges()); err != nil {
		return fmt.Errorf("save graph edges failed: %w", err)
	}
	snap := changes.NewSnapshot(changes.CurrentSnapshotID(), changes.CurrentSnapshotID())
	snap.Assets = e.registry.All()
	snap.Services = e.svcRegistry.All()
	snap.Findings = e.findReg.All()
	snap.Nodes = e.builder.Graph().Nodes()
	snap.Edges = e.builder.Graph().Edges()
	if err := st.SaveSnapshot(snap); err != nil {
		return fmt.Errorf("save snapshot failed: %w", err)
	}
	return nil
}

func (e *Engine) Load(path string) error {
	st, err := storage.Open(path)
	if err != nil {
		return fmt.Errorf("open store failed: %w", err)
	}
	defer st.Close()
	assets, err := st.LoadAssets()
	if err != nil {
		return fmt.Errorf("load assets failed: %w", err)
	}
	for _, a := range assets {
		e.registry.Add(a)
	}
	svcs, err := st.LoadServices()
	if err != nil {
		return fmt.Errorf("load services failed: %w", err)
	}
	for _, s := range svcs {
		e.svcRegistry.Add(s)
	}
	finds, err := st.LoadFindings()
	if err != nil {
		return fmt.Errorf("load findings failed: %w", err)
	}
	for _, f := range finds {
		e.findReg.Add(f)
	}
	nodes, err := st.LoadGraphNodes()
	if err != nil {
		return fmt.Errorf("load graph nodes failed: %w", err)
	}
	edges, err := st.LoadGraphEdges()
	if err != nil {
		return fmt.Errorf("load graph edges failed: %w", err)
	}
	e.builder.PopulateFromLoaded(e.registry.All(), e.svcRegistry.All(), e.findReg.All(), nodes, edges)
	return nil
}

func (e *Engine) Snapshot(path string) (*changes.Snapshot, error) {
	id := changes.CurrentSnapshotID()
	snap := changes.NewSnapshot(id, id)
	snap.Assets = e.registry.All()
	snap.Services = e.svcRegistry.All()
	snap.Findings = e.findReg.All()
	snap.Nodes = e.builder.Graph().Nodes()
	snap.Edges = e.builder.Graph().Edges()
	if path != "" {
		st, err := storage.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open store failed: %w", err)
		}
		defer st.Close()
		if err := st.SaveSnapshot(snap); err != nil {
			return nil, fmt.Errorf("save snapshot failed: %w", err)
		}
	}
	return snap, nil
}

func (e *Engine) LoadPreviousSnapshot(path string) (*changes.Snapshot, error) {
	st, err := storage.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open store failed: %w", err)
	}
	defer st.Close()
	return st.LoadPreviousSnapshot()
}
