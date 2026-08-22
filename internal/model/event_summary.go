package model

import "sort"

type EventSummary struct {
	FeederID      string   `json:"feeder_id"`
	Count         int      `json:"count"`
	MinFaultAmp   float64  `json:"min_fault_amp"`
	MaxFaultAmp   float64  `json:"max_fault_amp"`
	AverageAction float64  `json:"average_action_ms"`
	Relays        []string `json:"relays"`
}

func SummarizeEvents(feederID string, events []FaultEvent) EventSummary {
	summary := EventSummary{FeederID: feederID, Count: len(events)}
	if len(events) == 0 {
		return summary
	}
	set := map[string]bool{}
	for index, event := range events {
		if index == 0 || event.FaultAmp < summary.MinFaultAmp {
			summary.MinFaultAmp = event.FaultAmp
		}
		if event.FaultAmp > summary.MaxFaultAmp {
			summary.MaxFaultAmp = event.FaultAmp
		}
		summary.AverageAction += event.ActionMS
		if event.RelayID != "" {
			set[event.RelayID] = true
		}
	}
	summary.AverageAction /= float64(len(events))
	for relay := range set {
		summary.Relays = append(summary.Relays, relay)
	}
	sort.Strings(summary.Relays)
	return summary
}
