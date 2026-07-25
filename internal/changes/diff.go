package changes

import (
	"time"

	"github.com/networkscope/netscope/pkg/models"
)

type Change struct {
	Kind    ChangeKind
	Type    string
	ID      string
	Summary string
}

type ChangeResult struct {
	Changes []Change
	Summary map[ChangeKind]int
}

func Diff(before, after *Snapshot) *ChangeResult {
	res := &ChangeResult{Changes: make([]Change, 0), Summary: map[ChangeKind]int{}}

	assetsBefore := indexByID(before.Assets)
	assetsAfter := indexByID(after.Assets)
	for id, a := range assetsAfter {
		if _, ok := assetsBefore[id]; !ok {
			res.Changes = append(res.Changes, Change{Kind: Added, Type: "asset", ID: a.ID, Summary: string(a.Type)})
		}
	}
	for id, a := range assetsBefore {
		if _, ok := assetsAfter[id]; !ok {
			res.Changes = append(res.Changes, Change{Kind: Removed, Type: "asset", ID: a.ID, Summary: string(a.Type)})
		}
	}

	svcsBefore := indexByServiceID(before.Services)
	svcsAfter := indexByServiceID(after.Services)
	for id, s := range svcsAfter {
		if _, ok := svcsBefore[id]; !ok {
			res.Changes = append(res.Changes, Change{Kind: Added, Type: "service", ID: s.ID, Summary: s.Transport + "/" + s.Protocol})
		}
	}
	for id, s := range svcsBefore {
		if _, ok := svcsAfter[id]; !ok {
			res.Changes = append(res.Changes, Change{Kind: Removed, Type: "service", ID: s.ID, Summary: s.Transport + "/" + s.Protocol})
		}
	}

	for _, c := range res.Changes {
		res.Summary[c.Kind]++
	}
	return res
}

func indexByID(as []*models.Asset) map[string]*models.Asset {
	out := make(map[string]*models.Asset, len(as))
	for _, a := range as {
		out[a.ID] = a
	}
	return out
}

func indexByServiceID(ss []*models.Service) map[string]*models.Service {
	out := make(map[string]*models.Service, len(ss))
	for _, s := range ss {
		out[s.ID] = s
	}
	return out
}

func NewSnapshot(id, ts string) *Snapshot {
	return &Snapshot{ID: id, Timestamp: ts}
}

func CurrentSnapshotID() string {
	return time.Now().UTC().Format(time.RFC3339)
}
