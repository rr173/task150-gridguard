package model

import (
	"fmt"
	"time"
)

type PlanStatus string

const (
	PlanDraft     PlanStatus = "draft"
	PlanValidated PlanStatus = "validated"
	PlanActive    PlanStatus = "active"
	PlanRejected  PlanStatus = "rejected"
	PlanRetired   PlanStatus = "retired"
)

type RelayRole string

const (
	RelayPrimary  RelayRole = "primary"
	RelayBackup   RelayRole = "backup"
	RelayBoundary RelayRole = "boundary"
)

type Feeder struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	NominalAmp float64   `json:"nominal_amp"`
	UpstreamID string    `json:"upstream_id,omitempty"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (f Feeder) Validate() error {
	if f.ID == "" {
		return FieldError("feeder.id", "不能为空")
	}
	if f.Name == "" {
		return FieldError("feeder.name", "不能为空")
	}
	if f.NominalAmp <= 0 {
		return FieldError("feeder.nominal_amp", "必须大于 0")
	}
	if f.UpstreamID == f.ID {
		return FieldError("feeder.upstream_id", "不能指向自身")
	}
	return nil
}

type ProtectionZone struct {
	ID              string    `json:"id"`
	FeederID        string    `json:"feeder_id"`
	Name            string    `json:"name"`
	Sequence        int       `json:"sequence"`
	MinFaultAmp     float64   `json:"min_fault_amp"`
	MaxFaultAmp     float64   `json:"max_fault_amp"`
	IsolationSwitch string    `json:"isolation_switch"`
	CreatedAt       time.Time `json:"created_at"`
}

func (z ProtectionZone) Validate() error {
	if z.ID == "" || z.FeederID == "" || z.Name == "" {
		return FieldError("zone", "id、feeder_id 和 name 必填")
	}
	if z.Sequence < 1 {
		return FieldError("zone.sequence", "必须从 1 开始")
	}
	if z.MinFaultAmp <= 0 || z.MaxFaultAmp <= 0 || z.MinFaultAmp > z.MaxFaultAmp {
		return FieldError("zone.fault_range", "故障电流范围无效")
	}
	if z.IsolationSwitch == "" {
		return FieldError("zone.isolation_switch", "不能为空")
	}
	return nil
}

type RelaySetting struct {
	ID         string    `json:"id"`
	FeederID   string    `json:"feeder_id"`
	ZoneID     string    `json:"zone_id"`
	Name       string    `json:"name"`
	Role       RelayRole `json:"role"`
	PickupAmp  float64   `json:"pickup_amp"`
	InstantAmp float64   `json:"instant_amp"`
	TimeDial   float64   `json:"time_dial"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
}

func (r RelaySetting) Validate() error {
	if r.ID == "" || r.FeederID == "" || r.ZoneID == "" || r.Name == "" {
		return FieldError("relay", "id、feeder_id、zone_id 和 name 必填")
	}
	if r.Role != RelayPrimary && r.Role != RelayBackup && r.Role != RelayBoundary {
		return FieldError("relay.role", "不是支持的角色")
	}
	if r.PickupAmp <= 0 || r.InstantAmp <= r.PickupAmp {
		return FieldError("relay.pickup_amp", "必须小于瞬时阈值且均大于 0")
	}
	if r.TimeDial <= 0 {
		return FieldError("relay.time_dial", "必须大于 0")
	}
	return nil
}

type CoordinationPlan struct {
	ID          string     `json:"id"`
	FeederID    string     `json:"feeder_id"`
	Version     int        `json:"version"`
	Status      PlanStatus `json:"status"`
	MarginMS    float64    `json:"margin_ms"`
	CreatedAt   time.Time  `json:"created_at"`
	ValidatedAt *time.Time `json:"validated_at,omitempty"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
}

func (p CoordinationPlan) Validate() error {
	if p.ID == "" || p.FeederID == "" {
		return FieldError("plan", "id 和 feeder_id 必填")
	}
	if p.Version < 1 {
		return FieldError("plan.version", "必须大于 0")
	}
	if p.MarginMS <= 0 {
		return FieldError("plan.margin_ms", "必须大于 0")
	}
	if !p.Status.Valid() {
		return FieldError("plan.status", "不是支持的状态")
	}
	return nil
}

func (s PlanStatus) Valid() bool {
	return s == PlanDraft || s == PlanValidated || s == PlanRejected || s == PlanRetired
}

type CoordinationPair struct {
	ID           string  `json:"id"`
	PlanID       string  `json:"plan_id"`
	ZoneID       string  `json:"zone_id"`
	PrimaryRelay string  `json:"primary_relay"`
	BackupRelay  string  `json:"backup_relay"`
	FaultAmp     float64 `json:"fault_amp"`
	PrimaryMS    float64 `json:"primary_ms"`
	BackupMS     float64 `json:"backup_ms"`
	MarginMS     float64 `json:"margin_ms"`
	Passes       bool    `json:"passes"`
}

func (p CoordinationPair) SelectivityMS() float64 { return p.BackupMS - p.PrimaryMS }

type Violation struct {
	ID        string    `json:"id"`
	PlanID    string    `json:"plan_id"`
	ZoneID    string    `json:"zone_id,omitempty"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	CreatedAt time.Time `json:"created_at"`
}

func (v Violation) Error() string { return fmt.Sprintf("%s: %s", v.Code, v.Message) }

type FaultEvent struct {
	ID          string    `json:"id"`
	FeederID    string    `json:"feeder_id"`
	ZoneID      string    `json:"zone_id"`
	FaultAmp    float64   `json:"fault_amp"`
	OccurredAt  time.Time `json:"occurred_at"`
	PlanVersion int       `json:"plan_version"`
	RelayID     string    `json:"relay_id,omitempty"`
	ActionMS    float64   `json:"action_ms,omitempty"`
}

func (e FaultEvent) Validate() error {
	if e.ID == "" || e.FeederID == "" || e.ZoneID == "" {
		return FieldError("event", "id、feeder_id 和 zone_id 必填")
	}
	if e.FaultAmp <= 0 {
		return FieldError("event.fault_amp", "必须大于 0")
	}
	return nil
}

type ActiveVersion struct {
	FeederID string    `json:"feeder_id"`
	PlanID   string    `json:"plan_id"`
	Version  int       `json:"version"`
	Changed  time.Time `json:"changed_at"`
}
