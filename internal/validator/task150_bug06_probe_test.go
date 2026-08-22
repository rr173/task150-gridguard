package validator_test

import (
	"testing"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/validator"
)

func TestBug06_BlockingViolationStopsPlan(t *testing.T) {
	if !validator.HasBlocking([]model.Violation{{Severity:"error"}}){t.Fatal("error violation must block activation")}
}
