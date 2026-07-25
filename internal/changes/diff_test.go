package changes

import (
	"testing"

	"github.com/networkscope/netscope/pkg/models"
)

func TestDiffAddedRemoved(t *testing.T) {
	before := &Snapshot{
		Assets: []*models.Asset{mustAsset("a1", models.AssetTypeIP)},
	}
	after := &Snapshot{
		Assets: []*models.Asset{
			mustAsset("a1", models.AssetTypeIP),
			mustAsset("a2", models.AssetTypeDomain),
		},
	}
	result := Diff(before, after)
	if result.Summary[Added] != 1 {
		t.Errorf("Added = %d, want 1", result.Summary[Added])
	}
	if result.Summary[Removed] != 0 {
		t.Errorf("Removed = %d, want 0", result.Summary[Removed])
	}
	found := false
	for _, c := range result.Changes {
		if c.Kind == Added && c.ID == "a2" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected added change for a2")
	}
}

func TestDiffServices(t *testing.T) {
	before := &Snapshot{
		Services: []*models.Service{mustSvc("s1", "a1", 22, "tcp")},
	}
	after := &Snapshot{}
	result := Diff(before, after)
	if result.Summary[Removed] != 1 {
		t.Errorf("Removed = %d, want 1", result.Summary[Removed])
	}
}

func TestDiffModifiedSummary(t *testing.T) {
	a1before, _ := models.NewAsset("a1", models.AssetTypeIP, "dns")
	a1after, _ := models.NewAsset("a1", models.AssetTypeIP, "cert")
	before := &Snapshot{Assets: []*models.Asset{a1before}}
	after := &Snapshot{Assets: []*models.Asset{a1after}}
	result := Diff(before, after)
	if len(result.Changes) != 0 {
		t.Errorf("changes = %d, want 0 (same object identity still matches)", len(result.Changes))
	}
}

func TestDiffMixed(t *testing.T) {
	before := &Snapshot{
		Assets: []*models.Asset{
			mustAsset("a1", models.AssetTypeIP),
			mustAsset("a2", models.AssetTypeDomain),
		},
		Services: []*models.Service{
			mustSvc("s1", "a1", 22, "tcp"),
		},
	}
	after := &Snapshot{
		Assets: []*models.Asset{
			mustAsset("a1", models.AssetTypeIP),
			mustAsset("a3", models.AssetTypeHostname),
		},
		Services: []*models.Service{
			mustSvc("s2", "a1", 80, "tcp"),
		},
	}
	result := Diff(before, after)
	if result.Summary[Added] != 2 {
		t.Errorf("Added = %d, want 2", result.Summary[Added])
	}
	if result.Summary[Removed] != 2 {
		t.Errorf("Removed = %d, want 2", result.Summary[Removed])
	}
}

func mustAsset(id string, t models.AssetType) *models.Asset {
	a, err := models.NewAsset(id, t, "scan")
	if err != nil {
		panic(err)
	}
	return a
}

func mustSvc(id, assetID string, port int, transport string) *models.Service {
	s, err := models.NewService(id, assetID, port, transport)
	if err != nil {
		panic(err)
	}
	return s
}
