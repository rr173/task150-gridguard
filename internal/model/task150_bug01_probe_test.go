package model_test

import (
	"testing"

	"github.com/rr173/task150-gridguard/internal/curve"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/planner"
)

func TestBug01_SelectivityMarginKeepsBackupDelay(t *testing.T) {
	pair := model.CoordinationPair{PrimaryMS: 100, BackupMS: 220, MarginMS: 80}
	if got := pair.SelectivityMS(); got != 120 {
		t.Fatalf("backup selectivity = %vms, want 120ms", got)
	}
	if !curve.MarginOK(pair.PrimaryMS, pair.BackupMS, pair.MarginMS) {
		t.Fatal("120ms backup delay should satisfy an 80ms selectivity margin")
	}
	if got := planner.SelectivityDeficit(pair); got != 0 {
		t.Fatalf("satisfied pair deficit = %vms, want 0", got)
	}
}
