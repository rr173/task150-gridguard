package validator

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"sort"
)

func Unique(input []model.Violation) []model.Violation {
	seen := map[string]model.Violation{}
	for _, item := range input {
		if item.ID == "" {
			item.ID = item.PlanID + ":" + item.ZoneID + ":" + item.Code
		}
		seen[item.ID] = item
	}
	result := make([]model.Violation, 0, len(seen))
	for _, item := range seen {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].ZoneID == result[j].ZoneID {
			return result[i].Code < result[j].Code
		}
		return result[i].ZoneID < result[j].ZoneID
	})
	return result
}

func HasBlocking(items []model.Violation) bool {
	for _, item := range items {
		if item.Severity == "error" {
			return true
		}
	}
	return false
}

func CountByCode(items []model.Violation) map[string]int {
	result := map[string]int{}
	for _, item := range items {
		result[item.Code]++
	}
	return result
}

func ForZone(items []model.Violation, zoneID string) []model.Violation {
	result := []model.Violation{}
	for _, item := range items {
		if item.ZoneID == zoneID {
			result = append(result, item)
		}
	}
	return result
}
