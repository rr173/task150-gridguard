package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rr173/task150-gridguard/internal/model"
	"time"
)

func (s *Store) NextVersion(ctx context.Context, feederID string) (int, error) {
	var version int
	err := s.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version_no),0)+1 FROM plans WHERE feeder_id=?`, feederID).Scan(&version)
	return version, err
}

func (s *Store) SavePlan(ctx context.Context, plan model.CoordinationPlan, pairs []model.CoordinationPair) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO plans(id,feeder_id,version_no,status,margin_ms,created_at) VALUES(?,?,?,?,?,?)`, plan.ID, plan.FeederID, plan.Version, plan.Status, plan.MarginMS, nowText(plan.CreatedAt))
	if err != nil {
		return databaseError(err)
	}
	for _, pair := range pairs {
		_, err = tx.ExecContext(ctx, `INSERT INTO coordination_pairs(id,plan_id,zone_id,primary_relay,backup_relay,fault_amp,primary_ms,backup_ms,margin_ms,passes) VALUES(?,?,?,?,?,?,?,?,?,?)`, pair.ID, pair.PlanID, pair.ZoneID, pair.PrimaryRelay, pair.BackupRelay, pair.FaultAmp, pair.PrimaryMS, pair.BackupMS, pair.MarginMS, pair.Passes)
		if err != nil {
			return databaseError(err)
		}
	}
	return tx.Commit()
}

func (s *Store) SaveValidation(ctx context.Context, planID string, status model.PlanStatus, violations []model.Violation) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := nowText(time.Now())
	if _, err = tx.ExecContext(ctx, `DELETE FROM violations WHERE plan_id=?`, planID); err != nil {
		return err
	}
	for _, v := range violations {
		if v.CreatedAt.IsZero() {
			v.CreatedAt = time.Now()
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO violations(id,plan_id,zone_id,code,message,severity,created_at) VALUES(?,?,?,?,?,?,?)`, v.ID, v.PlanID, v.ZoneID, v.Code, v.Message, v.Severity, nowText(v.CreatedAt))
		if err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `UPDATE plans SET status=?,validated_at=? WHERE id=?`, status, now, planID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) Plan(ctx context.Context, id string) (model.CoordinationPlan, error) {
	var p model.CoordinationPlan
	var created string
	var validated, activated sql.NullString
	err := s.db.QueryRowContext(ctx, `SELECT id,feeder_id,version_no,status,margin_ms,created_at,validated_at,activated_at FROM plans WHERE id=?`, id).Scan(&p.ID, &p.FeederID, &p.Version, &p.Status, &p.MarginMS, &created, &validated, &activated)
	if err != nil {
		return p, oneRowNotFound(err, "计划", id)
	}
	p.CreatedAt = parseTime(created)
	if validated.Valid {
		t := parseTime(validated.String)
		p.ValidatedAt = &t
	}
	if activated.Valid {
		t := parseTime(activated.String)
		p.ActivatedAt = &t
	}
	return p, nil
}

func (s *Store) ListPlans(ctx context.Context) ([]model.CoordinationPlan, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,feeder_id,version_no,status,margin_ms,created_at,validated_at,activated_at FROM plans ORDER BY feeder_id,version_no`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.CoordinationPlan{}
	for rows.Next() {
		var p model.CoordinationPlan
		var created string
		var validated, activated sql.NullString
		if err := rows.Scan(&p.ID, &p.FeederID, &p.Version, &p.Status, &p.MarginMS, &created, &validated, &activated); err != nil {
			return nil, err
		}
		p.CreatedAt = parseTime(created)
		if validated.Valid {
			t := parseTime(validated.String)
			p.ValidatedAt = &t
		}
		if activated.Valid {
			t := parseTime(activated.String)
			p.ActivatedAt = &t
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

func (s *Store) Pairs(ctx context.Context, planID string) ([]model.CoordinationPair, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,plan_id,zone_id,primary_relay,backup_relay,fault_amp,primary_ms,backup_ms,margin_ms,passes FROM coordination_pairs WHERE plan_id=? ORDER BY zone_id`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.CoordinationPair{}
	for rows.Next() {
		var p model.CoordinationPair
		if err := rows.Scan(&p.ID, &p.PlanID, &p.ZoneID, &p.PrimaryRelay, &p.BackupRelay, &p.FaultAmp, &p.PrimaryMS, &p.BackupMS, &p.MarginMS, &p.Passes); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

func (s *Store) Violations(ctx context.Context, planID string) ([]model.Violation, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,plan_id,zone_id,code,message,severity,created_at FROM violations WHERE plan_id=? ORDER BY zone_id,code`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.Violation{}
	for rows.Next() {
		var v model.Violation
		var created string
		if err := rows.Scan(&v.ID, &v.PlanID, &v.ZoneID, &v.Code, &v.Message, &v.Severity, &created); err != nil {
			return nil, err
		}
		v.CreatedAt = parseTime(created)
		result = append(result, v)
	}
	return result, rows.Err()
}

func (s *Store) Activate(ctx context.Context, plan model.CoordinationPlan) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := nowText(time.Now())
	if _, err = tx.ExecContext(ctx, `UPDATE plans SET status=? WHERE feeder_id=? AND status=?`, model.PlanRetired, plan.FeederID, model.PlanActive); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE plans SET status=?,activated_at=? WHERE id=? AND status=?`, model.PlanActive, now, plan.ID, model.PlanValidated); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO active_versions(feeder_id,plan_id,version_no,changed_at) VALUES(?,?,?,?) ON CONFLICT(feeder_id) DO UPDATE SET plan_id=excluded.plan_id,version_no=excluded.version_no,changed_at=excluded.changed_at`, plan.FeederID, plan.ID, plan.Version, now); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) Active(ctx context.Context, feederID string) (model.ActiveVersion, error) {
	var a model.ActiveVersion
	var changed string
	err := s.db.QueryRowContext(ctx, `SELECT feeder_id,plan_id,version_no,changed_at FROM active_versions WHERE feeder_id=?`, feederID).Scan(&a.FeederID, &a.PlanID, &a.Version, &changed)
	if err != nil {
		return a, oneRowNotFound(err, "活动计划", feederID)
	}
	a.Changed = parseTime(changed)
	return a, nil
}

func (s *Store) SaveEvent(ctx context.Context, event model.FaultEvent) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO fault_events(id,feeder_id,zone_id,fault_amp,occurred_at,plan_version,relay_id,action_ms) VALUES(?,?,?,?,?,?,?,?)`, event.ID, event.FeederID, event.ZoneID, event.FaultAmp, nowText(event.OccurredAt), event.PlanVersion, event.RelayID, event.ActionMS)
	return err
}

func (s *Store) Stats(ctx context.Context) (model.SystemStats, error) {
	var x model.SystemStats
	queries := []struct {
		target *int
		sql    string
	}{{&x.Feeders, "SELECT COUNT(*) FROM feeders"}, {&x.Zones, "SELECT COUNT(*) FROM zones"}, {&x.Relays, "SELECT COUNT(*) FROM relays"}, {&x.Plans, "SELECT COUNT(*) FROM plans"}, {&x.Active, "SELECT COUNT(*) FROM active_versions"}, {&x.Events, "SELECT COUNT(*) FROM fault_events"}}
	for _, q := range queries {
		if err := s.db.QueryRowContext(ctx, q.sql).Scan(q.target); err != nil {
			return x, fmt.Errorf("统计: %w", err)
		}
	}
	return x, nil
}
