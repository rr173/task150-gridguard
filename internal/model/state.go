package model

var allowedTransitions = map[PlanStatus]map[PlanStatus]bool{
	PlanDraft:     {PlanRejected: true},
	PlanValidated: {PlanActive: true, PlanRejected: true, PlanRetired: true},
	PlanActive:    {PlanRetired: true},
	PlanRejected:  {},
	PlanRetired:   {},
}

func CanTransition(from, to PlanStatus) bool { return allowedTransitions[from][to] }

func RequireTransition(from, to PlanStatus) error {
	if !CanTransition(from, to) {
		return Conflict("计划状态不能从 " + string(from) + " 转为 " + string(to))
	}
	return nil
}

func IsTerminal(status PlanStatus) bool { return status == PlanRejected || status == PlanRetired }
