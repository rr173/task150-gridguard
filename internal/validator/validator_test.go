package validator

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"testing"
)

func TestUniqueDropsDuplicateEvidence(t *testing.T) {
	items := Unique([]model.Violation{{ID: "same", Code: "a"}, {ID: "same", Code: "a"}})
	if len(items) != 1 {
		t.Fatalf("items=%v", items)
	}
}
func TestHasBlocking(t *testing.T) {
	if !HasBlocking([]model.Violation{{Severity: "error"}}) {
		t.Fatal("expected blocking")
	}
	if HasBlocking([]model.Violation{{Severity: "warning"}}) {
		t.Fatal("warning is not blocking")
	}
}
