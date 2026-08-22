package topology

import "github.com/rr173/task150-gridguard/internal/model"

func (i *Index) CheckRelay(relay model.RelaySetting) error {
	if err := relay.Validate(); err != nil {
		return err
	}
	if _, ok := i.feeders[relay.FeederID]; !ok {
		return model.NotFound("馈线", relay.FeederID)
	}
	zone, ok := i.zones[relay.ZoneID]
	if !ok {
		return model.NotFound("保护区", relay.ZoneID)
	}
	if relay.Role == model.RelayPrimary && zone.FeederID != relay.FeederID {
		return model.Conflict("主保护继电器必须位于保护区所属馈线")
	}
	if relay.Role == model.RelayBackup && !i.IsUpstream(relay.FeederID, zone.FeederID) {
		return model.Conflict("后备继电器必须在保护区上游路径")
	}
	return nil
}

func (i *Index) CheckPlanFeeder(feederID string) error {
	feeder, ok := i.feeders[feederID]
	if !ok {
		return model.NotFound("馈线", feederID)
	}
	if feeder.Enabled {
		return model.Conflict("馈线已停用")
	}
	if len(i.zonesByFeeder[feederID]) == 0 {
		return model.Conflict("馈线没有保护区")
	}
	return nil
}

func (i *Index) ReachableZones(feederID string) ([]model.ProtectionZone, error) {
	if err := i.CheckPlanFeeder(feederID); err != nil {
		return nil, err
	}
	return i.Zones(feederID), nil
}
