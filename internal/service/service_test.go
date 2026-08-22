package service

import (
	"context"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/store"
	"path/filepath"
	"testing"
)

func TestPlanPersistsAndRecovers(t *testing.T) {
	ctx := context.Background()
	database, err := store.Open(filepath.Join(t.TempDir(), "grid.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	app, err := New(database)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = app.CreateFeeder(ctx, model.CreateFeederRequest{ID: "f", Name: "f", NominalAmp: 500}); err != nil {
		t.Fatal(err)
	}
	if _, err = app.CreateZone(ctx, model.CreateZoneRequest{ID: "z", FeederID: "f", Name: "z", Sequence: 1, MinFaultAmp: 100, MaxFaultAmp: 800, IsolationSwitch: "q"}); err != nil {
		t.Fatal(err)
	}
	if _, err = app.CreateRelay(ctx, model.CreateRelayRequest{ID: "p", FeederID: "f", ZoneID: "z", Name: "p", Role: model.RelayPrimary, PickupAmp: 80, InstantAmp: 700, TimeDial: .1}); err != nil {
		t.Fatal(err)
	}
	if _, err = app.CreateRelay(ctx, model.CreateRelayRequest{ID: "b", FeederID: "f", ZoneID: "z", Name: "b", Role: model.RelayBackup, PickupAmp: 40, InstantAmp: 1000, TimeDial: .8}); err != nil {
		t.Fatal(err)
	}
	plan, err := app.CreatePlan(ctx, model.CreatePlanRequest{ID: "plan", FeederID: "f", MarginMS: 80})
	if err != nil {
		t.Fatal(err)
	}
	checked, err := app.ValidatePlan(ctx, plan.Plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if checked.Plan.Status != model.PlanValidated {
		t.Fatalf("violations=%v", checked.Violations)
	}
	if _, err = app.ActivatePlan(ctx, plan.Plan.ID); err != nil {
		t.Fatal(err)
	}
	if err = app.Recover(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err = app.EvaluateEvent(ctx, model.EvaluateEventRequest{ID: "event", FeederID: "f", ZoneID: "z", FaultAmp: 450}); err != nil {
		t.Fatal(err)
	}
}
