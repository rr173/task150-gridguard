package store

import (
	"context"
	"database/sql"
	"github.com/rr173/task150-gridguard/internal/model"
)

func (s *Store) Events(ctx context.Context, feederID string, limit int) ([]model.FaultEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	query := `SELECT id,feeder_id,zone_id,fault_amp,occurred_at,plan_version,relay_id,action_ms FROM fault_events`
	args := []any{}
	if feederID != "" {
		query += ` WHERE feeder_id=?`
		args = append(args, feederID)
	}
	query += ` ORDER BY occurred_at DESC LIMIT ?`
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.FaultEvent{}
	for rows.Next() {
		item, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) Event(ctx context.Context, id string) (model.FaultEvent, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id,feeder_id,zone_id,fault_amp,occurred_at,plan_version,relay_id,action_ms FROM fault_events WHERE id=?`, id)
	item, err := scanEvent(row)
	if err != nil {
		return item, oneRowNotFound(err, "故障事件", id)
	}
	return item, nil
}

type eventScanner interface{ Scan(...any) error }

func scanEvent(row eventScanner) (model.FaultEvent, error) {
	var item model.FaultEvent
	var occurred string
	err := row.Scan(&item.ID, &item.FeederID, &item.ZoneID, &item.FaultAmp, &occurred, &item.PlanVersion, &item.RelayID, &item.ActionMS)
	if err != nil {
		return item, err
	}
	item.OccurredAt = parseTime(occurred)
	return item, nil
}

func (s *Store) EventCount(ctx context.Context, feederID string) (int, error) {
	var count int
	var err error
	if feederID == "" {
		err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM fault_events`).Scan(&count)
	} else {
		err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM fault_events WHERE feeder_id=?`, feederID).Scan(&count)
	}
	return count, err
}

func (s *Store) DeleteEvent(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM fault_events WHERE id=?`, id)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return model.NotFound("故障事件", id)
	}
	return nil
}

func nullableEvent(row *sql.Row) (model.FaultEvent, bool, error) {
	item, err := scanEvent(row)
	if err == sql.ErrNoRows {
		return model.FaultEvent{}, false, nil
	}
	return item, err == nil, err
}
