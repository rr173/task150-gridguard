package service

import (
	"context"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/topology"
	"time"
)

func (s *Service) CreateFeeder(ctx context.Context, request model.CreateFeederRequest) (model.Feeder, error) {
	f := model.Feeder{ID: request.ID, Name: request.Name, NominalAmp: request.NominalAmp, UpstreamID: request.UpstreamID, Enabled: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := f.Validate(); err != nil {
		return f, err
	}
	current, err := s.store.ListFeeders(ctx)
	if err != nil {
		return f, err
	}
	zones, err := s.store.ListZones(ctx)
	if err != nil {
		return f, err
	}
	if _, err := topology.New(append(current, f), zones); err != nil {
		return f, err
	}
	if err := s.store.PutFeeder(ctx, f); err != nil {
		return f, err
	}
	return f, s.Recover(ctx)
}

func (s *Service) CreateZone(ctx context.Context, request model.CreateZoneRequest) (model.ProtectionZone, error) {
	z := model.ProtectionZone{ID: request.ID, FeederID: request.FeederID, Name: request.Name, Sequence: request.Sequence, MinFaultAmp: request.MinFaultAmp, MaxFaultAmp: request.MaxFaultAmp, IsolationSwitch: request.IsolationSwitch, CreatedAt: time.Now()}
	if err := z.Validate(); err != nil {
		return z, err
	}
	index := s.currentIndex()
	if index == nil {
		return z, model.Conflict("服务尚未恢复")
	}
	if _, ok := index.Feeder(z.FeederID); !ok {
		return z, model.NotFound("馈线", z.FeederID)
	}
	if err := s.store.PutZone(ctx, z); err != nil {
		return z, err
	}
	return z, s.Recover(ctx)
}

func (s *Service) CreateRelay(ctx context.Context, request model.CreateRelayRequest) (model.RelaySetting, error) {
	r := model.RelaySetting{ID: request.ID, FeederID: request.FeederID, ZoneID: request.ZoneID, Name: request.Name, Role: request.Role, PickupAmp: request.PickupAmp, InstantAmp: request.InstantAmp, TimeDial: request.TimeDial, Enabled: true, CreatedAt: time.Now()}
	index := s.currentIndex()
	if index == nil {
		return r, model.Conflict("服务尚未恢复")
	}
	if err := index.CheckRelay(r); err != nil {
		return r, err
	}
	if err := s.store.PutRelay(ctx, r); err != nil {
		return r, err
	}
	return r, nil
}
