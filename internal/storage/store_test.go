package storage

import (
	"testing"

	"github.com/networkscope/netscope/pkg/models"
)

func TestStoreRoundtrip(t *testing.T) {
	type testCase struct {
		name      string
		fixture   func() *Store
		validator func(t *testing.T, s *Store)
	}
	
	tests := []testCase{
		{
			name: "in_memory_db",
			fixture: func() *Store {
				db, err := Open(":memory:")
				if err != nil {
					t.Fatalf("open failed: %v", err)
				}
				return db
			},
			validator: func(t *testing.T, s *Store) {
				a, err := models.NewAsset("a1", models.AssetTypeIP, "scan")
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				if err := s.SaveAsset(a); err != nil {
					t.Fatalf("save asset failed: %v", err)
				}
				loaded, err := s.LoadAssets()
				if err != nil {
					t.Fatalf("load assets failed: %v", err)
				}
				if len(loaded) != 1 {
					t.Errorf("loaded assets = %d, want 1", len(loaded))
				}
				if loaded[0].ID != "a1" {
					t.Errorf("ID = %q, want a1", loaded[0].ID)
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.fixture()
			defer store.Close()
			tt.validator(t, store)
		})
	}
}

func TestStoreClear(t *testing.T) {
	s, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer s.Close()
	
	a, err := models.NewAsset("a1", models.AssetTypeIP, "scan")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if err := s.SaveAsset(a); err != nil {
		t.Fatalf("save asset failed: %v", err)
	}
	
	if err := s.Clear(); err != nil {
		t.Fatalf("clear failed: %v", err)
	}
	
	loaded, err := s.LoadAssets()
	if err != nil {
		t.Fatalf("load after clear failed: %v", err)
	}
	if len(loaded) != 0 {
		t.Errorf("expected empty store after clear, got %d assets", len(loaded))
	}
}
