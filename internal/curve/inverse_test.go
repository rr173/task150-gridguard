package curve

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"testing"
)

func testRelay() model.RelaySetting {
	return model.RelaySetting{ID: "r", FeederID: "f", ZoneID: "z", Name: "r", Role: model.RelayPrimary, PickupAmp: 100, InstantAmp: 1000, TimeDial: .2, Enabled: true}
}
func TestInverseTimeDropsAtHigherCurrent(t *testing.T) {
	relay := testRelay()
	low, err := OperateMS(relay, 200)
	if err != nil {
		t.Fatal(err)
	}
	high, err := OperateMS(relay, 500)
	if err != nil {
		t.Fatal(err)
	}
	if high >= low {
		t.Fatalf("high=%v low=%v", high, low)
	}
}
func TestInstantOperation(t *testing.T) {
	value, err := OperateMS(testRelay(), 1000)
	if err != nil || value != 30 {
		t.Fatalf("%v %v", value, err)
	}
}
