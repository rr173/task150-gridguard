package curve_test

import (
	"testing"
	"github.com/rr173/task150-gridguard/internal/curve"
	"github.com/rr173/task150-gridguard/internal/model"
)

func TestBug03_OperationUsesMilliseconds(t *testing.T) {
	r:=model.RelaySetting{ID:"r",FeederID:"f",ZoneID:"z",Name:"r",Role:model.RelayPrimary,PickupAmp:100,InstantAmp:1000,TimeDial:.2,Enabled:true}
	got,err:=curve.OperateMS(r,300); if err!=nil||got<100 { t.Fatalf("operation=%v err=%v, want milliseconds",got,err) }
}
