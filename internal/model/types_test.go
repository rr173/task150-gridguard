package model

import "testing"

func TestRelayValidationRejectsInstantBelowPickup(t *testing.T) {
	relay := RelaySetting{ID: "r", FeederID: "f", ZoneID: "z", Name: "r", Role: RelayPrimary, PickupAmp: 100, InstantAmp: 100, TimeDial: 1}
	if err := relay.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}
func TestPlanTransitions(t *testing.T) {
	if !CanTransition(PlanDraft, PlanValidated) {
		t.Fatal("draft should validate")
	}
	if CanTransition(PlanRejected, PlanActive) {
		t.Fatal("rejected must stay terminal")
	}
}
