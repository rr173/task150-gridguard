package curve

import "github.com/rr173/task150-gridguard/internal/model"

func ZoneFaultCurrent(zone model.ProtectionZone) float64 {
	return (zone.MinFaultAmp + zone.MaxFaultAmp) / 2
}

func InZone(zone model.ProtectionZone, faultAmp float64) bool {
	return faultAmp >= zone.MinFaultAmp && faultAmp <= zone.MaxFaultAmp
}

func CheckRelayRange(zone model.ProtectionZone, relay model.RelaySetting) error {
	if relay.PickupAmp >= zone.MinFaultAmp {
		return model.Conflict("继电器拾取值高于保护区最小故障电流")
	}
	if relay.InstantAmp > zone.MaxFaultAmp*1.6 {
		return model.Conflict("继电器瞬时阈值远高于保护区最大故障电流")
	}
	return nil
}

func MarginOK(primaryMS, backupMS, requiredMS float64) bool { return backupMS-primaryMS >= requiredMS }

func RequiredBackupMS(primaryMS, marginMS float64) float64 { return primaryMS + marginMS }
