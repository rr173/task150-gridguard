package planner

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/topology"
	"testing"
)

func TestBuildPairs(t *testing.T) {
	feeders := []model.Feeder{{ID: "f", Name: "f", NominalAmp: 100, Enabled: true}}
	zones := []model.ProtectionZone{{ID: "z", FeederID: "f", Name: "z", Sequence: 1, MinFaultAmp: 100, MaxFaultAmp: 600, IsolationSwitch: "q"}}
	index, err := topology.New(feeders, zones)
	if err != nil {
		t.Fatal(err)
	}
	relays := []model.RelaySetting{{ID: "p", FeederID: "f", ZoneID: "z", Name: "p", Role: model.RelayPrimary, PickupAmp: 50, InstantAmp: 700, TimeDial: .1, Enabled: true}, {ID: "b", FeederID: "f", ZoneID: "z", Name: "b", Role: model.RelayBackup, PickupAmp: 30, InstantAmp: 900, TimeDial: .8, Enabled: true}}
	pairs, violations := New(index).Build(model.CoordinationPlan{ID: "plan", FeederID: "f", Version: 1, Status: model.PlanDraft, MarginMS: 80}, relays)
	if len(violations) != 0 || len(pairs) != 1 {
		t.Fatalf("pairs=%v violations=%v", pairs, violations)
	}
}

// TestBuildPairsUpstreamBackup reproduces the reported defect: when the
// backup relay sits on the upstream feeder (the correct place for backup
// protection), the previous reversed direction check dropped it and the
// planner reported backup_missing. With the fix the upstream backup path
// participates in coordination, and a backup on a downstream feeder is
// rejected instead.
func TestBuildPairsUpstreamBackup(t *testing.T) {
	feeders := []model.Feeder{
		{ID: "src", Name: "src", NominalAmp: 100, Enabled: true},
		{ID: "f", Name: "f", NominalAmp: 100, UpstreamID: "src", Enabled: true},
		{ID: "g", Name: "g", NominalAmp: 100, UpstreamID: "f", Enabled: true},
	}
	zones := []model.ProtectionZone{{ID: "z", FeederID: "f", Name: "z", Sequence: 1, MinFaultAmp: 100, MaxFaultAmp: 600, IsolationSwitch: "q"}}
	index, err := topology.New(feeders, zones)
	if err != nil {
		t.Fatal(err)
	}
	relays := []model.RelaySetting{
		{ID: "p", FeederID: "f", ZoneID: "z", Name: "p", Role: model.RelayPrimary, PickupAmp: 50, InstantAmp: 700, TimeDial: .1, Enabled: true},
		// valid backup relay located on the upstream feeder — must be kept
		{ID: "b", FeederID: "src", ZoneID: "z", Name: "b", Role: model.RelayBackup, PickupAmp: 30, InstantAmp: 900, TimeDial: .8, Enabled: true},
		// invalid backup relay on a downstream feeder — must be dropped
		{ID: "d", FeederID: "g", ZoneID: "z", Name: "d", Role: model.RelayBackup, PickupAmp: 20, InstantAmp: 950, TimeDial: .05, Enabled: true},
	}
	pairs, violations := New(index).Build(model.CoordinationPlan{ID: "plan", FeederID: "f", Version: 1, Status: model.PlanDraft, MarginMS: 80}, relays)
	if len(pairs) != 1 {
		t.Fatalf("expected one pair, got pairs=%v violations=%v", pairs, violations)
	}
	if pairs[0].BackupRelay != "b" {
		t.Fatalf("expected upstream backup b, got %s", pairs[0].BackupRelay)
	}
	for _, v := range violations {
		if v.Code == "backup_missing" {
			t.Fatalf("upstream backup was dropped: %v", v)
		}
	}
}
