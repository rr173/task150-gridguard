package topology_test

import (
	"testing"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/topology"
)

func TestBug08_EnabledFeederAcceptsPlan(t *testing.T) {
	i,err:=topology.New([]model.Feeder{{ID:"f",Name:"f",NominalAmp:1,Enabled:true}},[]model.ProtectionZone{{ID:"z",FeederID:"f",Name:"z",Sequence:1,MinFaultAmp:1,MaxFaultAmp:2,IsolationSwitch:"q"}});if err!=nil{t.Fatal(err)};if err=i.CheckPlanFeeder("f");err!=nil{t.Fatalf("enabled feeder rejected: %v",err)}
}
