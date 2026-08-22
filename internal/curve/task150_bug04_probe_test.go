package curve_test

import (
	"testing"
	"github.com/rr173/task150-gridguard/internal/curve"
	"github.com/rr173/task150-gridguard/internal/model"
)

func TestBug04_FaultBoundaryIsInclusive(t *testing.T) {
	z:=model.ProtectionZone{MinFaultAmp:100,MaxFaultAmp:800}; if !curve.InZone(z,100)||!curve.InZone(z,800){t.Fatal("fault range boundaries must be accepted")}
}
