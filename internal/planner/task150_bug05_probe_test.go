package planner_test

import (
	"testing"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/planner"
)

func TestBug05_PassedPairCount(t *testing.T) {
	if got:=planner.PassedPairs([]model.CoordinationPair{{Passes:true},{Passes:true},{Passes:false}});got!=2{t.Fatalf("passed=%d, want 2",got)}
}
