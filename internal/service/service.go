package service

import (
	"context"
	"github.com/rr173/task150-gridguard/internal/model"
	"github.com/rr173/task150-gridguard/internal/store"
	"github.com/rr173/task150-gridguard/internal/topology"
	"sync"
)

type Service struct {
	store       *store.Store
	mu          sync.RWMutex
	index       *topology.Index
	feederLocks map[string]*sync.Mutex
}

func New(database *store.Store) (*Service, error) {
	s := &Service{store: database, feederLocks: map[string]*sync.Mutex{}}
	if err := s.Recover(context.Background()); err != nil {
		return nil, err
	}
	return s, nil
}
func (s *Service) Store() *store.Store { return s.store }

func (s *Service) Recover(ctx context.Context) error {
	feeders, err := s.store.ListFeeders(ctx)
	if err != nil {
		return err
	}
	zones, err := s.store.ListZones(ctx)
	if err != nil {
		return err
	}
	index, err := topology.New(feeders, zones)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.index = index
	for _, f := range feeders {
		if s.feederLocks[f.ID] == nil {
			s.feederLocks[f.ID] = &sync.Mutex{}
		}
	}
	s.mu.Unlock()
	return nil
}

func (s *Service) currentIndex() *topology.Index { s.mu.RLock(); defer s.mu.RUnlock(); return s.index }
func (s *Service) lockFeeder(id string) func() {
	s.mu.Lock()
	lock := s.feederLocks[id]
	if lock == nil {
		lock = &sync.Mutex{}
		s.feederLocks[id] = lock
	}
	s.mu.Unlock()
	lock.Lock()
	return lock.Unlock
}

func (s *Service) Feeders(ctx context.Context) ([]model.Feeder, error) {
	return s.store.ListFeeders(ctx)
}
func (s *Service) Zones(ctx context.Context) ([]model.ProtectionZone, error) {
	return s.store.ListZones(ctx)
}
func (s *Service) Relays(ctx context.Context) ([]model.RelaySetting, error) {
	return s.store.ListRelays(ctx)
}
func (s *Service) Plans(ctx context.Context) ([]model.CoordinationPlan, error) {
	return s.store.ListPlans(ctx)
}
func (s *Service) Stats(ctx context.Context) (model.SystemStats, error) { return s.store.Stats(ctx) }
func (s *Service) Events(ctx context.Context, feederID string, limit int) ([]model.FaultEvent, error) {
	return s.store.Events(ctx, feederID, limit)
}
func (s *Service) EventSummary(ctx context.Context, feederID string) (model.EventSummary, error) {
	events, err := s.store.Events(ctx, feederID, 1)
	if err != nil {
		return model.EventSummary{}, err
	}
	return model.SummarizeEvents(feederID, events), nil
}
