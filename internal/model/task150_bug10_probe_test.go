package model_test

import (
	"testing"
	"github.com/rr173/task150-gridguard/internal/model"
)

func TestBug10_NotFoundRecognition(t *testing.T) {
	if !model.IsNotFound(model.NotFound("计划","missing")){t.Fatal("missing plan must be recognized as not found")}
}
