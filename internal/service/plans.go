package service

import (
	"context"
	"github.com/rr173/task150-gridguard/internal/curve"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/planner"
	"github.com/rr173/task150-gridguard/internal/validator"
	"time"
)

func (s *Service) CreatePlan(ctx context.Context, request model.CreatePlanRequest) (model.PlanSummary, error) {
	unlock := s.lockFeeder(request.FeederID)
	defer unlock()
	index := s.currentIndex()
	if index == nil {
		return model.PlanSummary{}, model.Conflict("服务尚未恢复")
	}
	if err := index.CheckPlanFeeder(request.FeederID); err != nil {
		return model.PlanSummary{}, err
	}
	version, err := s.store.NextVersion(ctx, request.FeederID)
	if err != nil {
		return model.PlanSummary{}, err
	}
	plan := model.CoordinationPlan{ID: request.ID, FeederID: request.FeederID, Version: version, Status: model.PlanDraft, MarginMS: request.MarginMS, CreatedAt: time.Now()}
	if err := plan.Validate(); err != nil {
		return model.PlanSummary{}, err
	}
	relays, err := s.store.ListRelays(ctx)
	if err != nil {
		return model.PlanSummary{}, err
	}
	pairs, initial := planner.New(index).Build(plan, relays)
	if err := s.store.SavePlan(ctx, plan, pairs); err != nil {
		return model.PlanSummary{}, err
	}
	return model.PlanSummary{Plan: plan, Pairs: planner.SortPairs(pairs), Violations: initial}, nil
}

func (s *Service) ValidatePlan(ctx context.Context, planID string) (model.PlanSummary, error) {
	plan, err := s.store.Plan(ctx, planID)
	if err != nil {
		return model.PlanSummary{}, err
	}
	unlock := s.lockFeeder(plan.FeederID)
	defer unlock()
	index := s.currentIndex()
	relays, err := s.store.ListRelays(ctx)
	if err != nil {
		return model.PlanSummary{}, err
	}
	pairs, initial := planner.New(index).Build(plan, relays)
	violations := validator.New(index).Check(plan, relays, pairs, initial)
	status := model.PlanValidated
	if validator.HasBlocking(violations) {
		status = model.PlanRejected
	}
	if err := s.store.SaveValidation(ctx, planID, status, violations); err != nil {
		return model.PlanSummary{}, err
	}
	plan.Status = status
	now := time.Now()
	plan.ValidatedAt = &now
	return model.PlanSummary{Plan: plan, Pairs: planner.SortPairs(pairs), Violations: violations}, nil
}

func (s *Service) Plan(ctx context.Context, id string) (model.PlanSummary, error) {
	plan, err := s.store.Plan(ctx, id)
	if err != nil {
		return model.PlanSummary{}, err
	}
	pairs, err := s.store.Pairs(ctx, id)
	if err != nil {
		return model.PlanSummary{}, err
	}
	violations, err := s.store.Violations(ctx, id)
	if err != nil {
		return model.PlanSummary{}, err
	}
	return model.PlanSummary{Plan: plan, Pairs: pairs, Violations: violations}, nil
}

func (s *Service) ActivatePlan(ctx context.Context, id string) (model.ActiveVersion, error) {
	plan, err := s.store.Plan(ctx, id)
	if err != nil {
		return model.ActiveVersion{}, err
	}
	unlock := s.lockFeeder(plan.FeederID)
	defer unlock()
	if plan.Status != model.PlanValidated {
		return model.ActiveVersion{}, model.Conflict("只有已校验计划可以激活")
	}
	violations, err := s.store.Violations(ctx, id)
	if err != nil {
		return model.ActiveVersion{}, err
	}
	if validator.HasBlocking(violations) {
		return model.ActiveVersion{}, model.Conflict("计划存在阻断违例")
	}
	if err := s.store.Activate(ctx, plan); err != nil {
		return model.ActiveVersion{}, err
	}
	return s.store.Active(ctx, plan.FeederID)
}

func (s *Service) EvaluateEvent(ctx context.Context, request model.EvaluateEventRequest) (model.FaultEvent, error) {
	event := model.FaultEvent{ID: request.ID, FeederID: request.FeederID, ZoneID: request.ZoneID, FaultAmp: request.FaultAmp, OccurredAt: time.Now()}
	if err := event.Validate(); err != nil {
		return event, err
	}
	active, err := s.store.Active(ctx, event.FeederID)
	if err != nil {
		return event, err
	}
	zoneIndex := s.currentIndex()
	zone, ok := zoneIndex.Zone(event.ZoneID)
	if !ok || zone.FeederID != event.FeederID {
		return event, model.NotFound("馈线内保护区", event.ZoneID)
	}
	if !curve.InZone(zone, event.FaultAmp) {
		return event, model.FieldError("fault_amp", "不在保护区故障电流范围")
	}
	pairs, err := s.store.Pairs(ctx, active.PlanID)
	if err != nil {
		return event, err
	}
	pair, ok := planner.PairByZone(pairs, event.ZoneID)
	if !ok {
		return event, model.Conflict("活动计划没有该保护区的协调对")
	}
	event.PlanVersion = active.Version
	event.RelayID = pair.PrimaryRelay
	event.ActionMS = pair.PrimaryMS
	if err := s.store.SaveEvent(ctx, event); err != nil {
		return event, err
	}
	return event, nil
}
