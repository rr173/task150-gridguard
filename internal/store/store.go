package store

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"time"
)

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("打开 sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)
	s := &Store{db: db}
	if err := s.Migrate(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }
func (s *Store) DB() *sql.DB  { return s.db }

func (s *Store) Migrate(ctx context.Context) error {
	statements := []string{
		`PRAGMA foreign_keys = ON`,
		`CREATE TABLE IF NOT EXISTS feeders (id TEXT PRIMARY KEY, name TEXT NOT NULL, nominal_amp REAL NOT NULL, upstream_id TEXT NOT NULL DEFAULT '', enabled INTEGER NOT NULL, created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS zones (id TEXT PRIMARY KEY, feeder_id TEXT NOT NULL, name TEXT NOT NULL, sequence_no INTEGER NOT NULL, min_fault_amp REAL NOT NULL, max_fault_amp REAL NOT NULL, isolation_switch TEXT NOT NULL, created_at TEXT NOT NULL, UNIQUE(feeder_id, sequence_no), FOREIGN KEY(feeder_id) REFERENCES feeders(id))`,
		`CREATE TABLE IF NOT EXISTS relays (id TEXT PRIMARY KEY, feeder_id TEXT NOT NULL, zone_id TEXT NOT NULL, name TEXT NOT NULL, role TEXT NOT NULL, pickup_amp REAL NOT NULL, instant_amp REAL NOT NULL, time_dial REAL NOT NULL, enabled INTEGER NOT NULL, created_at TEXT NOT NULL, FOREIGN KEY(feeder_id) REFERENCES feeders(id), FOREIGN KEY(zone_id) REFERENCES zones(id))`,
		`CREATE TABLE IF NOT EXISTS plans (id TEXT PRIMARY KEY, feeder_id TEXT NOT NULL, version_no INTEGER NOT NULL, status TEXT NOT NULL, margin_ms REAL NOT NULL, created_at TEXT NOT NULL, validated_at TEXT, activated_at TEXT, UNIQUE(feeder_id, version_no), FOREIGN KEY(feeder_id) REFERENCES feeders(id))`,
		`CREATE TABLE IF NOT EXISTS coordination_pairs (id TEXT PRIMARY KEY, plan_id TEXT NOT NULL, zone_id TEXT NOT NULL, primary_relay TEXT NOT NULL, backup_relay TEXT NOT NULL, fault_amp REAL NOT NULL, primary_ms REAL NOT NULL, backup_ms REAL NOT NULL, margin_ms REAL NOT NULL, passes INTEGER NOT NULL, FOREIGN KEY(plan_id) REFERENCES plans(id))`,
		`CREATE TABLE IF NOT EXISTS violations (id TEXT PRIMARY KEY, plan_id TEXT NOT NULL, zone_id TEXT NOT NULL DEFAULT '', code TEXT NOT NULL, message TEXT NOT NULL, severity TEXT NOT NULL, created_at TEXT NOT NULL, FOREIGN KEY(plan_id) REFERENCES plans(id))`,
		`CREATE TABLE IF NOT EXISTS active_versions (feeder_id TEXT PRIMARY KEY, plan_id TEXT NOT NULL, version_no INTEGER NOT NULL, changed_at TEXT NOT NULL, FOREIGN KEY(plan_id) REFERENCES plans(id))`,
		`CREATE TABLE IF NOT EXISTS fault_events (id TEXT PRIMARY KEY, feeder_id TEXT NOT NULL, zone_id TEXT NOT NULL, fault_amp REAL NOT NULL, occurred_at TEXT NOT NULL, plan_version INTEGER NOT NULL, relay_id TEXT NOT NULL DEFAULT '', action_ms REAL NOT NULL DEFAULT 0)`,
		`CREATE INDEX IF NOT EXISTS idx_zones_feeder ON zones(feeder_id, sequence_no)`,
		`CREATE INDEX IF NOT EXISTS idx_relays_zone ON relays(zone_id)`,
		`CREATE INDEX IF NOT EXISTS idx_plans_feeder ON plans(feeder_id, version_no)`,
		`CREATE INDEX IF NOT EXISTS idx_violations_plan ON violations(plan_id)`,
	}
	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("迁移数据库: %w", err)
		}
	}
	return nil
}

func nowText(t time.Time) string { return t.UTC().Format(time.RFC3339Nano) }
func parseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}
	}
	return t
}
