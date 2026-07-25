package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/networkscope/netscope/internal/changes"
	"github.com/networkscope/netscope/pkg/models"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening store: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	return &Store{db: db}, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		`PRAGMA journal_mode=WAL;`,
		`PRAGMA foreign_keys=ON;`,
		`CREATE TABLE IF NOT EXISTS assets (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			source TEXT NOT NULL,
			evidence TEXT DEFAULT '',
			first_seen TEXT NOT NULL,
			last_seen TEXT NOT NULL,
			metadata TEXT DEFAULT '{}'
		);`,
		`CREATE TABLE IF NOT EXISTS services (
			id TEXT PRIMARY KEY,
			asset_id TEXT NOT NULL,
			port INTEGER NOT NULL CHECK(port BETWEEN 1 AND 65535),
			transport TEXT NOT NULL,
			protocol TEXT DEFAULT '',
			name TEXT DEFAULT '',
			software TEXT DEFAULT '',
			version TEXT DEFAULT '',
			confidence REAL NOT NULL DEFAULT 1.0,
			first_seen TEXT NOT NULL,
			last_seen TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS findings (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			severity TEXT NOT NULL,
			affected_asset TEXT NOT NULL,
			evidence TEXT DEFAULT '',
			description TEXT DEFAULT '',
			recommendation TEXT DEFAULT '',
			status TEXT NOT NULL DEFAULT 'open',
			confidence REAL NOT NULL DEFAULT 1.0,
			first_seen TEXT NOT NULL,
			last_seen TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS graph_nodes (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			label TEXT DEFAULT ''
		);`,
		`CREATE TABLE IF NOT EXISTS graph_edges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source TEXT NOT NULL,
			target TEXT NOT NULL,
			type TEXT NOT NULL,
			label TEXT DEFAULT ''
		);`,
		`CREATE TABLE IF NOT EXISTS snapshots (
			id TEXT PRIMARY KEY,
			timestamp TEXT NOT NULL,
			payload TEXT NOT NULL
		);`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) SaveAsset(a *models.Asset) error {
	if a == nil {
		return nil
	}
	_, err := s.db.Exec(`
		INSERT INTO assets (id, type, source, evidence, first_seen, last_seen, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET last_seen=excluded.last_seen, metadata=excluded.metadata`,
		a.ID, string(a.Type), a.Source, a.Evidence, a.FirstSeen.Format(time.RFC3339), a.LastSeen.Format(time.RFC3339), marshalMap(a.Metadata))
	return err
}

func (s *Store) SaveService(svc *models.Service) error {
	if svc == nil {
		return nil
	}
	_, err := s.db.Exec(`
		INSERT INTO services (id, asset_id, port, transport, protocol, name, software, version, confidence, first_seen, last_seen)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET last_seen=excluded.last_seen`,
		svc.ID, svc.AssetID, svc.Port, svc.Transport, svc.Protocol, svc.Name, svc.Software, svc.Version, svc.Confidence,
		svc.FirstSeen.Format(time.RFC3339), svc.LastSeen.Format(time.RFC3339))
	return err
}

func (s *Store) SaveFinding(f *models.Finding) error {
	if f == nil {
		return nil
	}
	_, err := s.db.Exec(`
		INSERT INTO findings (id, title, severity, affected_asset, evidence, description, recommendation, status, confidence, first_seen, last_seen)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET last_seen=excluded.last_seen, status=excluded.status`,
		f.ID, f.Title, string(f.Severity), f.AffectedAsset, f.Evidence, f.Description, f.Recommendation, string(f.Status), f.Confidence,
		f.FirstSeen.Format(time.RFC3339), f.LastSeen.Format(time.RFC3339))
	return err
}

func (s *Store) SaveGraphNodes(nodes []*models.Node) error {
	if len(nodes) == 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT INTO graph_nodes (id, type, label) VALUES (?, ?, ?) ON CONFLICT(id) DO NOTHING`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, n := range nodes {
		if _, err := stmt.Exec(n.ID, n.Type, n.Label); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) SaveGraphEdges(edges []models.Edge) error {
	if len(edges) == 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT INTO graph_edges (source, target, type, label) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, e := range edges {
		if _, err := stmt.Exec(e.Source, e.Target, e.Type, e.Label); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) LoadAssets() ([]*models.Asset, error) {
	rows, err := s.db.Query(`SELECT id, type, source, evidence, first_seen, last_seen, metadata FROM assets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Asset
	for rows.Next() {
		a := &models.Asset{Metadata: make(map[string]string)}
		var fs, ls, meta string
		if err := rows.Scan(&a.ID, (*string)(&a.Type), &a.Source, &a.Evidence, &fs, &ls, &meta); err != nil {
			return nil, err
		}
		if a.FirstSeen, err = time.Parse(time.RFC3339, fs); err != nil {
			return nil, err
		}
		if a.LastSeen, err = time.Parse(time.RFC3339, ls); err != nil {
			return nil, err
		}
		for k, v := range unmarshalMap(meta) {
			a.Metadata[k] = v
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) Clear() error {
	_, err := s.db.Exec(`DELETE FROM graph_edges; DELETE FROM graph_nodes; DELETE FROM findings; DELETE FROM services; DELETE FROM assets;`)
	return err
}

func marshalMap(m map[string]string) string {
	out := "{"
	first := true
	for k, v := range m {
		if !first {
			out += ","
		}
		out += fmt.Sprintf("%q:%q", k, v)
		first = false
	}
	out += "}"
	return out
}

func unmarshalMap(s string) map[string]string {
	if s == "" || s == "{}" {
		return make(map[string]string)
	}
	// minimal parser for "{key:val,...}"
	m := make(map[string]string)
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return m
	}
	s = s[1 : len(s)-1]
	key := ""
	val := ""
	readingKey := true
	escaped := false
	for _, ch := range s {
		switch {
		case escaped:
			key += string(ch)
			escaped = false
		case ch == '\\':
			escaped = true
		case readingKey && ch == ':':
			readingKey = false
		case !readingKey && ch == ',':
			m[key] = val
			key = ""
			val = ""
			readingKey = true
		default:
			if readingKey {
				key = key + string(ch)
			} else {
				val = val + string(ch)
			}
		}
	}
	if key != "" {
		m[key] = val
	}
	return m
}

func (s *Store) SaveSnapshot(snap *changes.Snapshot) error {
	if snap == nil {
		return nil
	}
	payload, err := json.Marshal(snap)
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}
	_, err = s.db.Exec(`INSERT INTO snapshots (id, timestamp, payload) VALUES (?, ?, ?)`,
		snap.ID, snap.Timestamp, string(payload))
	return err
}

func (s *Store) LoadSnapshot(id string) (*changes.Snapshot, error) {
	row := s.db.QueryRow(`SELECT payload FROM snapshots WHERE id = ?`, id)
	var payload string
	if err := row.Scan(&payload); err != nil {
		return nil, err
	}
	var snap changes.Snapshot
	if err := json.Unmarshal([]byte(payload), &snap); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return &snap, nil
}

func (s *Store) LoadLatestSnapshot() (*changes.Snapshot, error) {
	row := s.db.QueryRow(`SELECT payload FROM snapshots ORDER BY timestamp DESC LIMIT 1`)
	var payload string
	if err := row.Scan(&payload); err != nil {
		return nil, err
	}
	var snap changes.Snapshot
	if err := json.Unmarshal([]byte(payload), &snap); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return &snap, nil
}

func (s *Store) LoadPreviousSnapshot() (*changes.Snapshot, error) {
	row := s.db.QueryRow(`SELECT payload FROM snapshots ORDER BY timestamp DESC LIMIT 1 OFFSET 1`)
	var payload string
	if err := row.Scan(&payload); err != nil {
		return nil, err
	}
	var snap changes.Snapshot
	if err := json.Unmarshal([]byte(payload), &snap); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return &snap, nil
}

func (s *Store) LoadServices() ([]*models.Service, error) {
	rows, err := s.db.Query(`SELECT id, asset_id, port, transport, protocol, name, software, version, confidence, first_seen, last_seen FROM services`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Service
	for rows.Next() {
		var svc models.Service
		var fs, ls string
		if err := rows.Scan(&svc.ID, &svc.AssetID, &svc.Port, &svc.Transport, &svc.Protocol, &svc.Name, &svc.Software, &svc.Version, &svc.Confidence, &fs, &ls); err != nil {
			return nil, err
		}
		if svc.FirstSeen, err = time.Parse(time.RFC3339, fs); err != nil {
			return nil, err
		}
		if svc.LastSeen, err = time.Parse(time.RFC3339, ls); err != nil {
			return nil, err
		}
		out = append(out, &svc)
	}
	return out, rows.Err()
}

func (s *Store) LoadFindings() ([]*models.Finding, error) {
	rows, err := s.db.Query(`SELECT id, title, severity, affected_asset, evidence, description, recommendation, status, confidence, first_seen, last_seen FROM findings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Finding
	for rows.Next() {
		var f models.Finding
		var sev, status, fs, ls string
		if err := rows.Scan(&f.ID, &f.Title, &sev, &f.AffectedAsset, &f.Evidence, &f.Description, &f.Recommendation, &status, &f.Confidence, &fs, &ls); err != nil {
			return nil, err
		}
		f.Severity = models.Severity(sev)
		f.Status = models.FindingStatus(status)
		if f.FirstSeen, err = time.Parse(time.RFC3339, fs); err != nil {
			return nil, err
		}
		if f.LastSeen, err = time.Parse(time.RFC3339, ls); err != nil {
			return nil, err
		}
		out = append(out, &f)
	}
	return out, rows.Err()
}

func (s *Store) LoadGraphNodes() ([]*models.Node, error) {
	rows, err := s.db.Query(`SELECT id, type, label FROM graph_nodes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Node
	for rows.Next() {
		var n models.Node
		if err := rows.Scan(&n.ID, &n.Type, &n.Label); err != nil {
			return nil, err
		}
		out = append(out, &n)
	}
	return out, rows.Err()
}

func (s *Store) LoadGraphEdges() ([]models.Edge, error) {
	rows, err := s.db.Query(`SELECT source, target, type, label FROM graph_edges`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Edge
	for rows.Next() {
		var e models.Edge
		if err := rows.Scan(&e.Source, &e.Target, &e.Type, &e.Label); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

