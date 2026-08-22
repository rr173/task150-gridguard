package topology

import "github.com/rr173/task150-gridguard/internal/model"

func (i *Index) UpstreamPath(feederID string) ([]model.Feeder, error) {
	if _, ok := i.feeders[feederID]; !ok {
		return nil, model.NotFound("馈线", feederID)
	}
	path := []model.Feeder{}
	for current := feederID; current != ""; {
		feeder := i.feeders[current]
		path = append(path, feeder)
		current = feeder.UpstreamID
	}
	return path, nil
}

// IsUpstream reports whether candidateID lies on the upstream path of
// feederID, i.e. the chain reached by following UpstreamID upward from
// feederID. The path is inclusive: a feeder is on its own upstream path,
// so IsUpstream(f, f) is true.
func (i *Index) IsUpstream(candidateID, feederID string) bool {
	for current := feederID; current != ""; {
		if current == candidateID {
			return true
		}
		item, ok := i.feeders[current]
		if !ok {
			return false
		}
		current = item.UpstreamID
	}
	return false
}

func (i *Index) ZonePath(zoneID string) ([]model.ProtectionZone, error) {
	zone, ok := i.zones[zoneID]
	if !ok {
		return nil, model.NotFound("保护区", zoneID)
	}
	feeders, err := i.UpstreamPath(zone.FeederID)
	if err != nil {
		return nil, err
	}
	path := []model.ProtectionZone{}
	for _, feeder := range feeders {
		path = append(path, i.zonesByFeeder[feeder.ID]...)
	}
	return path, nil
}

func (i *Index) ContainsZone(feederID, zoneID string) bool {
	zone, ok := i.zones[zoneID]
	return ok && zone.FeederID == feederID
}
