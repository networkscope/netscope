package models

import (
	"errors"
)

type Node struct {
	ID    string
	Type  string
	Label string
}

type Edge struct {
	Source string
	Target string
	Type   string
	Label  string
}

type Graph struct {
	nodes map[string]*Node
	edges []Edge
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*Node),
		edges: make([]Edge, 0),
	}
}

func (g *Graph) AddNode(n *Node) error {
	if n == nil || n.ID == "" {
		return errors.New("node ID cannot be empty")
	}
	g.nodes[n.ID] = n
	return nil
}

func (g *Graph) AddEdge(source, target, edgeType string) error {
	if source == "" || target == "" || edgeType == "" {
		return errors.New("edge source, target, and type cannot be empty")
	}
	g.edges = append(g.edges, Edge{
		Source: source,
		Target: target,
		Type:   edgeType,
	})
	return nil
}

func (g *Graph) Nodes() []*Node {
	out := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		out = append(out, n)
	}
	return out
}

func (g *Graph) Edges() []Edge {
	return append([]Edge(nil), g.edges...)
}
