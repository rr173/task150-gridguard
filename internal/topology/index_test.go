package topology

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"testing"
)

func TestUpstreamPath(t *testing.T) {
	index, err := New([]model.Feeder{{ID: "root", Name: "root", NominalAmp: 1, Enabled: true}, {ID: "child", Name: "child", NominalAmp: 1, UpstreamID: "root", Enabled: true}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	path, err := index.UpstreamPath("child")
	if err != nil {
		t.Fatal(err)
	}
	if len(path) != 2 || path[1].ID != "root" {
		t.Fatalf("path=%v", path)
	}
}
func TestRejectsFeederCycle(t *testing.T) {
	_, err := New([]model.Feeder{{ID: "a", Name: "a", NominalAmp: 1, UpstreamID: "b", Enabled: true}, {ID: "b", Name: "b", NominalAmp: 1, UpstreamID: "a", Enabled: true}}, nil)
	if err == nil {
		t.Fatal("expected cycle")
	}
}
