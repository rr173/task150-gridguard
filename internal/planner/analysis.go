package planner

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"sort"
)

func Assess(plan model.CoordinationPlan, pairs []model.CoordinationPair, violations []model.Violation) model.FeederAssessment {
	assessment := model.FeederAssessment{FeederID: plan.FeederID, PlanID: plan.ID, Version: plan.Version, PlanStatus: plan.Status, ZoneCount: len(pairs), WorstMarginMS: 0}
	byZone := map[string][]model.Violation{}
	for _, item := range violations {
		byZone[item.ZoneID] = append(byZone[item.ZoneID], item)
		if item.Severity == "error" {
			assessment.BlockingViolations++
		}
	}
	for _, pair := range pairs {
		margin := pair.SelectivityMS()
		if assessment.WorstMarginMS == 0 || margin < assessment.WorstMarginMS {
			assessment.WorstMarginMS = margin
		}
		entry := model.ZoneAssessment{ZoneID: pair.ZoneID, PrimaryRelay: pair.PrimaryRelay, BackupRelay: pair.BackupRelay, FaultAmp: pair.FaultAmp, PrimaryMS: pair.PrimaryMS, BackupMS: pair.BackupMS, RequiredMarginMS: pair.MarginMS, ActualMarginMS: margin, ViolationCount: len(byZone[pair.ZoneID])}
		if pair.Passes && entry.ViolationCount == 0 {
			entry.Status = "pass"
			entry.Recommendation = "维持当前整定并持续监测故障事件"
			assessment.PassingPairs++
		} else {
			entry.Status = "blocked"
			entry.MarginDeficitMS = SelectivityDeficit(pair)
			entry.Recommendation = recommend(entry, byZone[pair.ZoneID])
		}
		assessment.Assessments = append(assessment.Assessments, entry)
	}
	for zoneID, items := range byZone {
		if _, ok := PairByZone(pairs, zoneID); ok {
			continue
		}
		assessment.Assessments = append(assessment.Assessments, model.ZoneAssessment{ZoneID: zoneID, Status: "blocked", ViolationCount: len(items), Recommendation: recommend(model.ZoneAssessment{}, items)})
	}
	sort.Slice(assessment.Assessments, func(i, j int) bool { return assessment.Assessments[i].ZoneID < assessment.Assessments[j].ZoneID })
	return assessment
}

func recommend(entry model.ZoneAssessment, items []model.Violation) string {
	for _, item := range items {
		switch item.Code {
		case "primary_missing":
			return "为该保护区补充可动作的主保护继电器"
		case "backup_missing":
			return "在上游馈线配置可动作的后备保护继电器"
		case "selectivity_shortfall", "margin_insufficient":
			return "调整主保护与后备保护的时间拨盘以恢复选择性裕度"
		case "primary_range", "backup_range":
			return "检查拾取值和瞬时阈值是否覆盖保护区故障电流范围"
		case "relay_topology":
			return "修正继电器所在馈线、保护区或角色的拓扑关系"
		}
	}
	if entry.MarginDeficitMS > 0 {
		return "提高后备保护延时或降低主保护延时"
	}
	return "检查保护区配置和继电器可用状态"
}

func Compare(current, candidate []model.CoordinationPair, currentPlan, candidatePlan string) model.PlanComparison {
	result := model.PlanComparison{CurrentPlanID: currentPlan, CandidatePlanID: candidatePlan}
	old := map[string]model.CoordinationPair{}
	for _, pair := range current {
		old[pair.ZoneID] = pair
	}
	for _, next := range candidate {
		prior, ok := old[next.ZoneID]
		if !ok {
			continue
		}
		result.SharedZones++
		delta := next.SelectivityMS() - prior.SelectivityMS()
		result.NetMarginChangeMS += delta
		if delta > 0 {
			result.ImprovedZones++
		}
		if delta < 0 {
			result.RegressedZones++
		}
	}
	return result
}

func RiskRank(items []model.ZoneAssessment) []model.ZoneAssessment {
	result := append([]model.ZoneAssessment(nil), items...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Status != result[j].Status {
			return result[i].Status == "blocked"
		}
		if result[i].MarginDeficitMS == result[j].MarginDeficitMS {
			return result[i].ZoneID < result[j].ZoneID
		}
		return result[i].MarginDeficitMS > result[j].MarginDeficitMS
	})
	return result
}
