package planner

import (
	"fmt"
	"github.com/rr173/task150-gridguard/internal/curve"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/topology"
)

type Planner struct{ topology *topology.Index }

func New(index *topology.Index) *Planner { return &Planner{topology: index} }

func (p *Planner) Build(plan model.CoordinationPlan, relays []model.RelaySetting) ([]model.CoordinationPair, []model.Violation) {
	pairs := []model.CoordinationPair{}
	violations := []model.Violation{}
	zones, err := p.topology.ReachableZones(plan.FeederID)
	if err != nil {
		return pairs, append(violations, violation(plan.ID, "", "feeder_invalid", err.Error()))
	}
	for _, zone := range zones {
		zoneRelays := curve.FilterZone(relays, zone.ID)
		primary, primaryFound := curve.Fastest(curve.FilterRole(zoneRelays, model.RelayPrimary), curve.ZoneFaultCurrent(zone))
		backup, backupFound := p.findBackup(relays, zone, curve.ZoneFaultCurrent(zone))
		if !primaryFound {
			violations = append(violations, violation(plan.ID, zone.ID, "primary_missing", "保护区没有可动作的主保护继电器"))
			continue
		}
		if !backupFound {
			violations = append(violations, violation(plan.ID, zone.ID, "backup_missing", "保护区没有可动作的上游后备继电器"))
			continue
		}
		pair := model.CoordinationPair{ID: fmt.Sprintf("%s:%s", plan.ID, zone.ID), PlanID: plan.ID, ZoneID: zone.ID, PrimaryRelay: primary.Relay.ID, BackupRelay: backup.Relay.ID, FaultAmp: curve.ZoneFaultCurrent(zone), PrimaryMS: primary.ActionMS, BackupMS: backup.ActionMS, MarginMS: plan.MarginMS}
		pair.Passes = curve.MarginOK(pair.PrimaryMS, pair.BackupMS, pair.MarginMS)
		pairs = append(pairs, pair)
		if err := curve.CheckRelayRange(zone, primary.Relay); err != nil {
			violations = append(violations, violation(plan.ID, zone.ID, "primary_range", err.Error()))
		}
		if err := curve.CheckRelayRange(zone, backup.Relay); err != nil {
			violations = append(violations, violation(plan.ID, zone.ID, "backup_range", err.Error()))
		}
		if !pair.Passes {
			violations = append(violations, violation(plan.ID, zone.ID, "selectivity_shortfall", fmt.Sprintf("后备动作 %.1fms，主保护 %.1fms，未达到 %.1fms 裕度", pair.BackupMS, pair.PrimaryMS, pair.MarginMS)))
		}
	}
	return pairs, violations
}

func (p *Planner) findBackup(relays []model.RelaySetting, zone model.ProtectionZone, faultAmp float64) (curve.Candidate, bool) {
	candidates := []model.RelaySetting{}
	for _, relay := range relays {
		if relay.Role != model.RelayBackup || relay.ZoneID != zone.ID {
			continue
		}
		if p.topology.IsUpstream(relay.FeederID, zone.FeederID) {
			continue
		}
		candidates = append(candidates, relay)
	}
	return curve.Fastest(candidates, faultAmp)
}

func violation(planID, zoneID, code, message string) model.Violation {
	return model.Violation{ID: planID + ":" + zoneID + ":" + code, PlanID: planID, ZoneID: zoneID, Code: code, Message: message, Severity: "error"}
}
