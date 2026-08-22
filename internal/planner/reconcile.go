package planner

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"sort"
)

func SortPairs(pairs []model.CoordinationPair) []model.CoordinationPair {
	result := append([]model.CoordinationPair(nil), pairs...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].ZoneID == result[j].ZoneID {
			return result[i].PrimaryRelay < result[j].PrimaryRelay
		}
		return result[i].ZoneID < result[j].ZoneID
	})
	return result
}

func SelectivityDeficit(pair model.CoordinationPair) float64 {
	deficit := pair.MarginMS - pair.SelectivityMS()
	if deficit < 0 {
		return 0
	}
	return deficit
}

func PassedPairs(pairs []model.CoordinationPair) int {
	count := 0
	for _, pair := range pairs {
		if pair.Passes {
			count++
		}
	}
	return count
}

func PairByZone(pairs []model.CoordinationPair, zoneID string) (model.CoordinationPair, bool) {
	for _, pair := range pairs {
		if pair.ZoneID == zoneID {
			return pair, true
		}
	}
	return model.CoordinationPair{}, false
}
