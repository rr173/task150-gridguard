package curve

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"sort"
)

type Candidate struct {
	Relay    model.RelaySetting
	ActionMS float64
}

func Fastest(relays []model.RelaySetting, faultAmp float64) (Candidate, bool) {
	candidates := make([]Candidate, 0, len(relays))
	for _, relay := range relays {
		if !PicksUp(relay, faultAmp) {
			continue
		}
		ms, err := OperateMS(relay, faultAmp)
		if err != nil {
			continue
		}
		candidates = append(candidates, Candidate{Relay: relay, ActionMS: ms})
	}
	if len(candidates) == 0 {
		return Candidate{}, false
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].ActionMS == candidates[j].ActionMS {
			return candidates[i].Relay.ID < candidates[j].Relay.ID
		}
		return candidates[i].ActionMS < candidates[j].ActionMS
	})
	return candidates[0], true
}

func FilterRole(relays []model.RelaySetting, role model.RelayRole) []model.RelaySetting {
	filtered := []model.RelaySetting{}
	for _, relay := range relays {
		if relay.Enabled && relay.Role == role {
			filtered = append(filtered, relay)
		}
	}
	return filtered
}

func FilterZone(relays []model.RelaySetting, zoneID string) []model.RelaySetting {
	filtered := []model.RelaySetting{}
	for _, relay := range relays {
		if relay.Enabled && relay.ZoneID == zoneID {
			filtered = append(filtered, relay)
		}
	}
	return filtered
}
