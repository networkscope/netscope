package changes

import (
	"github.com/networkscope/netscope/pkg/models"
)

type Snapshot struct {
	ID        string
	Timestamp string
	Assets    []*models.Asset
	Services  []*models.Service
	Findings  []*models.Finding
	Nodes     []*models.Node
	Edges     []models.Edge
}

type ChangeKind string

const (
	Added    ChangeKind = "added"
	Modified ChangeKind = "modified"
	Removed  ChangeKind = "removed"
)
