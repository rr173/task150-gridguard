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

func TestIsUpstreamDirection(t *testing.T) {
	index, err := New([]model.Feeder{
		{ID: "root", Name: "root", NominalAmp: 1, Enabled: true},
		{ID: "mid", Name: "mid", NominalAmp: 1, UpstreamID: "root", Enabled: true},
		{ID: "child", Name: "child", NominalAmp: 1, UpstreamID: "mid", Enabled: true},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		candidate, feeder string
		want              bool
	}{
		// upstream ancestors are on the path
		{"root", "child", true},  // grandparent upstream of child
		{"mid", "child", true},   // direct parent upstream of child
		{"root", "mid", true},    // parent upstream of mid
		// same feeder is on its own upstream path
		{"child", "child", true},
		{"root", "root", true},
		// downstream / sideways feeders are NOT upstream
		{"child", "root", false}, // child is downstream of root
		{"child", "mid", false},  // child is downstream of mid
		{"mid", "root", false},   // mid is downstream of root
		// unknown feeder is never upstream
		{"missing", "child", false},
		{"root", "missing", false},
	}
	for _, c := range cases {
		if got := index.IsUpstream(c.candidate, c.feeder); got != c.want {
			t.Errorf("IsUpstream(%q,%q)=%v want %v", c.candidate, c.feeder, got, c.want)
		}
	}
}

