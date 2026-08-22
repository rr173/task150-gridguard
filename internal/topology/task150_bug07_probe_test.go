package topology_test

import (
	"testing"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/topology"
)

func TestBug07_UpstreamBackupPath(t *testing.T) {
	i,err:=topology.New([]model.Feeder{{ID:"root",Name:"root",NominalAmp:1,Enabled:true},{ID:"child",Name:"child",NominalAmp:1,UpstreamID:"root",Enabled:true}},nil);if err!=nil{t.Fatal(err)};if i.IsUpstream("child","root"){t.Fatal("downstream feeder must not be treated as upstream")}
}
