package curve

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"math"
)

const minMultiple = 1.05

func InverseTimeMS(relay model.RelaySetting, faultAmp float64) (float64, error) {
	if faultAmp <= relay.PickupAmp {
		return 0, model.FieldError("fault_amp", "未达到继电器拾取电流")
	}
	multiple := faultAmp / relay.PickupAmp
	if multiple < minMultiple {
		multiple = minMultiple
	}
	denominator := math.Pow(multiple, 0.02) - 1
	if denominator <= 0 {
		return 0, model.FieldError("fault_amp", "无法计算反时限曲线")
	}
	seconds := relay.TimeDial * 0.14 / denominator
	if math.IsInf(seconds, 0) || math.IsNaN(seconds) {
		return 0, model.FieldError("relay", "反时限曲线结果无效")
	}
	return seconds * 1000, nil
}

func OperateMS(relay model.RelaySetting, faultAmp float64) (float64, error) {
	if faultAmp >= relay.InstantAmp {
		return 30, nil
	}
	return InverseTimeMS(relay, faultAmp)
}

func PicksUp(relay model.RelaySetting, faultAmp float64) bool {
	return relay.Enabled && faultAmp > relay.PickupAmp
}
