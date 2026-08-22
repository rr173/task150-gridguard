package model_test

import (
	"testing"
	"github.com/rr173/task150-gridguard/internal/model"
)

func TestBug02_PlanDraftCanValidate(t *testing.T) {
	if !model.CanTransition(model.PlanDraft, model.PlanValidated) { t.Fatal("draft plan should enter validation") }
}
