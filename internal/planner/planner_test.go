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
