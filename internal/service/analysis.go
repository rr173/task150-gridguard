package service

import (
	"context"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/planner"
)

func (s *Service) AssessPlan(ctx context.Context, planID string) (model.FeederAssessment, error) {
	summary, err := s.Plan(ctx, planID)
	if err != nil {
		return model.FeederAssessment{}, err
	}
	return planner.Assess(summary.Plan, summary.Pairs, summary.Violations), nil
}

func (s *Service) ComparePlans(ctx context.Context, currentID, candidateID string) (model.PlanComparison, error) {
	current, err := s.Plan(ctx, currentID)
	if err != nil {
		return model.PlanComparison{}, err
	}
	candidate, err := s.Plan(ctx, candidateID)
	if err != nil {
		return model.PlanComparison{}, err
	}
	if current.Plan.FeederID != candidate.Plan.FeederID {
		return model.PlanComparison{}, model.Conflict("只能比较同一馈线的计划")
	}
	return planner.Compare(current.Pairs, candidate.Pairs, current.Plan.ID, candidate.Plan.ID), nil
}

func (s *Service) ActiveAssessment(ctx context.Context, feederID string) (model.FeederAssessment, error) {
	active, err := s.store.Active(ctx, feederID)
	if err != nil {
		return model.FeederAssessment{}, err
	}
	return s.AssessPlan(ctx, active.PlanID)
}
