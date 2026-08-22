package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rr173/task150-gridguard/internal/model"
)

func (s *Store) PutFeeder(ctx context.Context, f model.Feeder) error {
	if err := f.Validate(); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO feeders(id,name,nominal_amp,upstream_id,enabled,created_at,updated_at) VALUES(?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,nominal_amp=excluded.nominal_amp,upstream_id=excluded.upstream_id,enabled=excluded.enabled,updated_at=excluded.updated_at`, f.ID, f.Name, f.NominalAmp, f.UpstreamID, f.Enabled, nowText(f.CreatedAt), nowText(f.UpdatedAt))
	return err
}

func (s *Store) ListFeeders(ctx context.Context) ([]model.Feeder, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,name,nominal_amp,upstream_id,enabled,created_at,updated_at FROM feeders ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.Feeder{}
	for rows.Next() {
		var f model.Feeder
		var created, updated string
		if err := rows.Scan(&f.ID, &f.Name, &f.NominalAmp, &f.UpstreamID, &f.Enabled, &created, &updated); err != nil {
			return nil, err
		}
		f.CreatedAt = parseTime(created)
		f.UpdatedAt = parseTime(updated)
		result = append(result, f)
	}
	return result, rows.Err()
}

func (s *Store) PutZone(ctx context.Context, z model.ProtectionZone) error {
	if err := z.Validate(); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO zones(id,feeder_id,name,sequence_no,min_fault_amp,max_fault_amp,isolation_switch,created_at) VALUES(?,?,?,?,?,?,?,?)`, z.ID, z.FeederID, z.Name, z.Sequence, z.MinFaultAmp, z.MaxFaultAmp, z.IsolationSwitch, nowText(z.CreatedAt))
	return err
}

func (s *Store) ListZones(ctx context.Context) ([]model.ProtectionZone, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,feeder_id,name,sequence_no,min_fault_amp,max_fault_amp,isolation_switch,created_at FROM zones ORDER BY feeder_id,sequence_no`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.ProtectionZone{}
	for rows.Next() {
		var z model.ProtectionZone
		var created string
		if err := rows.Scan(&z.ID, &z.FeederID, &z.Name, &z.Sequence, &z.MinFaultAmp, &z.MaxFaultAmp, &z.IsolationSwitch, &created); err != nil {
			return nil, err
		}
		z.CreatedAt = parseTime(created)
		result = append(result, z)
	}
	return result, rows.Err()
}

func (s *Store) PutRelay(ctx context.Context, r model.RelaySetting) error {
	if err := r.Validate(); err != nil {
		return err
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO relays(id,feeder_id,zone_id,name,role,pickup_amp,instant_amp,time_dial,enabled,created_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, r.ID, r.FeederID, r.ZoneID, r.Name, r.Role, r.PickupAmp, r.InstantAmp, r.TimeDial, r.Enabled, nowText(r.CreatedAt))
	return err
}

func (s *Store) ListRelays(ctx context.Context) ([]model.RelaySetting, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,feeder_id,zone_id,name,role,pickup_amp,instant_amp,time_dial,enabled,created_at FROM relays ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.RelaySetting{}
	for rows.Next() {
		var r model.RelaySetting
		var created string
		if err := rows.Scan(&r.ID, &r.FeederID, &r.ZoneID, &r.Name, &r.Role, &r.PickupAmp, &r.InstantAmp, &r.TimeDial, &r.Enabled, &created); err != nil {
			return nil, err
		}
		r.CreatedAt = parseTime(created)
		result = append(result, r)
	}
	return result, rows.Err()
}

func oneRowNotFound(err error, kind, id string) error {
	if err == sql.ErrNoRows {
		return model.NotFound(kind, id)
	}
	return err
}
func databaseError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("数据库操作失败: %w", err)
}
