package model

type CreateFeederRequest struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	NominalAmp float64 `json:"nominal_amp"`
	UpstreamID string  `json:"upstream_id"`
}
type CreateZoneRequest struct {
	ID              string  `json:"id"`
	FeederID        string  `json:"feeder_id"`
	Name            string  `json:"name"`
	Sequence        int     `json:"sequence"`
	MinFaultAmp     float64 `json:"min_fault_amp"`
	MaxFaultAmp     float64 `json:"max_fault_amp"`
	IsolationSwitch string  `json:"isolation_switch"`
}
type CreateRelayRequest struct {
	ID         string    `json:"id"`
	FeederID   string    `json:"feeder_id"`
	ZoneID     string    `json:"zone_id"`
	Name       string    `json:"name"`
	Role       RelayRole `json:"role"`
	PickupAmp  float64   `json:"pickup_amp"`
	InstantAmp float64   `json:"instant_amp"`
	TimeDial   float64   `json:"time_dial"`
}
type CreatePlanRequest struct {
	ID       string  `json:"id"`
	FeederID string  `json:"feeder_id"`
	MarginMS float64 `json:"margin_ms"`
}
type EvaluateEventRequest struct {
	ID       string  `json:"id"`
	FeederID string  `json:"feeder_id"`
	ZoneID   string  `json:"zone_id"`
	FaultAmp float64 `json:"fault_amp"`
}

type PlanSummary struct {
	Plan       CoordinationPlan   `json:"plan"`
	Pairs      []CoordinationPair `json:"pairs"`
	Violations []Violation        `json:"violations"`
}
type SystemStats struct {
	Feeders int `json:"feeders"`
	Zones   int `json:"zones"`
	Relays  int `json:"relays"`
	Plans   int `json:"plans"`
	Active  int `json:"active"`
	Events  int `json:"events"`
}
