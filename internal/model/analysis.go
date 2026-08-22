package model

type ZoneAssessment struct {
	ZoneID           string  `json:"zone_id"`
	PrimaryRelay     string  `json:"primary_relay,omitempty"`
	BackupRelay      string  `json:"backup_relay,omitempty"`
	FaultAmp         float64 `json:"fault_amp,omitempty"`
	PrimaryMS        float64 `json:"primary_ms,omitempty"`
	BackupMS         float64 `json:"backup_ms,omitempty"`
	RequiredMarginMS float64 `json:"required_margin_ms"`
	ActualMarginMS   float64 `json:"actual_margin_ms"`
	MarginDeficitMS  float64 `json:"margin_deficit_ms"`
	Status           string  `json:"status"`
	ViolationCount   int     `json:"violation_count"`
	Recommendation   string  `json:"recommendation"`
}

type FeederAssessment struct {
	FeederID           string           `json:"feeder_id"`
	PlanID             string           `json:"plan_id"`
	Version            int              `json:"version"`
	PlanStatus         PlanStatus       `json:"plan_status"`
	ZoneCount          int              `json:"zone_count"`
	PassingPairs       int              `json:"passing_pairs"`
	BlockingViolations int              `json:"blocking_violations"`
	WorstMarginMS      float64          `json:"worst_margin_ms"`
	Assessments        []ZoneAssessment `json:"assessments"`
}

type PlanComparison struct {
	CurrentPlanID     string  `json:"current_plan_id"`
	CandidatePlanID   string  `json:"candidate_plan_id"`
	SharedZones       int     `json:"shared_zones"`
	ImprovedZones     int     `json:"improved_zones"`
	RegressedZones    int     `json:"regressed_zones"`
	NetMarginChangeMS float64 `json:"net_margin_change_ms"`
}

func (a FeederAssessment) Healthy() bool {
	return a.BlockingViolations == 0 && a.ZoneCount > 0 && a.PassingPairs == a.ZoneCount
}
func (a ZoneAssessment) NeedsAction() bool { return a.Status != "pass" }
