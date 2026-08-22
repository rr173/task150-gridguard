package validator

import (
	"fmt"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/topology"
)

type Validator struct{ topology *topology.Index }

func New(index *topology.Index) *Validator { return &Validator{topology: index} }

func (v *Validator) Check(plan model.CoordinationPlan, relays []model.RelaySetting, pairs []model.CoordinationPair, initial []model.Violation) []model.Violation {
	violations := append([]model.Violation(nil), initial...)
	for _, relay := range relays {
		if !relay.Enabled {
			continue
		}
		if err := v.topology.CheckRelay(relay); err != nil {
			violations = append(violations, record(plan.ID, relay.ZoneID, "relay_topology", err.Error()))
		}
	}
	for _, pair := range pairs {
		if pair.PrimaryRelay == pair.BackupRelay {
			violations = append(violations, record(plan.ID, pair.ZoneID, "role_collision", "主保护与后备保护不能是同一个继电器"))
		}
		if pair.PrimaryMS <= 0 || pair.BackupMS <= 0 {
			violations = append(violations, record(plan.ID, pair.ZoneID, "action_time_invalid", "动作时间必须大于 0"))
		}
		if pair.BackupMS <= pair.PrimaryMS {
			violations = append(violations, record(plan.ID, pair.ZoneID, "backup_not_slower", "后备保护必须慢于主保护"))
		}
		if pair.SelectivityMS() < pair.MarginMS {
			violations = append(violations, record(plan.ID, pair.ZoneID, "margin_insufficient", fmt.Sprintf("选择性间隔 %.1fms 小于要求 %.1fms", pair.SelectivityMS(), pair.MarginMS)))
		}
	}
	return Unique(violations)
}

func record(planID, zoneID, code, message string) model.Violation {
	return model.Violation{ID: planID + ":" + zoneID + ":" + code, PlanID: planID, ZoneID: zoneID, Code: code, Message: message, Severity: "error"}
}
