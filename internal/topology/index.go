package topology

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"sort"
)

type Index struct {
	feeders       map[string]model.Feeder
	zones         map[string]model.ProtectionZone
	zonesByFeeder map[string][]model.ProtectionZone
}

func New(feeders []model.Feeder, zones []model.ProtectionZone) (*Index, error) {
	i := &Index{feeders: map[string]model.Feeder{}, zones: map[string]model.ProtectionZone{}, zonesByFeeder: map[string][]model.ProtectionZone{}}
	for _, feeder := range feeders {
		if err := feeder.Validate(); err != nil {
			return nil, err
		}
		if _, exists := i.feeders[feeder.ID]; exists {
			return nil, model.Conflict("重复馈线 " + feeder.ID)
		}
		i.feeders[feeder.ID] = feeder
	}
	for _, zone := range zones {
		if err := zone.Validate(); err != nil {
			return nil, err
		}
		if _, ok := i.feeders[zone.FeederID]; !ok {
			return nil, model.NotFound("馈线", zone.FeederID)
		}
		if _, exists := i.zones[zone.ID]; exists {
			return nil, model.Conflict("重复保护区 " + zone.ID)
		}
		i.zones[zone.ID] = zone
		i.zonesByFeeder[zone.FeederID] = append(i.zonesByFeeder[zone.FeederID], zone)
	}
	for feederID, items := range i.zonesByFeeder {
		sort.Slice(items, func(a, b int) bool { return items[a].Sequence < items[b].Sequence })
		for n := 1; n < len(items); n++ {
			if items[n-1].Sequence == items[n].Sequence {
				return nil, model.Conflict("馈线 " + feederID + " 的保护区顺序重复")
			}
		}
		i.zonesByFeeder[feederID] = items
	}
	if err := i.ValidateFeederChains(); err != nil {
		return nil, err
	}
	return i, nil
}

func (i *Index) Feeder(id string) (model.Feeder, bool) { value, ok := i.feeders[id]; return value, ok }
func (i *Index) Zone(id string) (model.ProtectionZone, bool) {
	value, ok := i.zones[id]
	return value, ok
}
func (i *Index) Zones(feederID string) []model.ProtectionZone {
	return append([]model.ProtectionZone(nil), i.zonesByFeeder[feederID]...)
}
func (i *Index) Feeders() []model.Feeder {
	values := make([]model.Feeder, 0, len(i.feeders))
	for _, value := range i.feeders {
		values = append(values, value)
	}
	sort.Slice(values, func(a, b int) bool { return values[a].ID < values[b].ID })
	return values
}

func (i *Index) ValidateFeederChains() error {
	for id := range i.feeders {
		seen := map[string]bool{}
		current := id
		for current != "" {
			if seen[current] {
				return model.Conflict("馈线拓扑存在环: " + id)
			}
			seen[current] = true
			feeder, ok := i.feeders[current]
			if !ok {
				return model.NotFound("上游馈线", current)
			}
			current = feeder.UpstreamID
		}
	}
	return nil
}
