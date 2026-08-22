package model_test

import (
	"testing"
	"github.com/rr173/task150-gridguard/internal/model"
)

func TestBug09_EventAverageAction(t *testing.T) {
	s:=model.SummarizeEvents("f",[]model.FaultEvent{{ActionMS:20},{ActionMS:40}});if s.AverageAction!=30{t.Fatalf("average=%v, want 30",s.AverageAction)}
}
