package httpapi

import (
	"github.com/rr173/task150-gridguard/internal/model"
	"net/http"
)

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	write(w, http.StatusOK, map[string]string{"status": "ok"})
}
func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	if _, err := s.service.Stats(r.Context()); err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, map[string]string{"status": "ready"})
}
func (s *Server) feeders(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	value, err := s.service.Feeders(r.Context())
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) createFeeder(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var req model.CreateFeederRequest
	if err := decode(r, &req); err != nil {
		fail(w, err)
		return
	}
	value, err := s.service.CreateFeeder(r.Context(), req)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusCreated, value)
}
func (s *Server) zones(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	value, err := s.service.Zones(r.Context())
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) createZone(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var req model.CreateZoneRequest
	if err := decode(r, &req); err != nil {
		fail(w, err)
		return
	}
	value, err := s.service.CreateZone(r.Context(), req)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusCreated, value)
}
func (s *Server) relays(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	value, err := s.service.Relays(r.Context())
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) createRelay(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var req model.CreateRelayRequest
	if err := decode(r, &req); err != nil {
		fail(w, err)
		return
	}
	value, err := s.service.CreateRelay(r.Context(), req)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusCreated, value)
}
func (s *Server) plans(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	value, err := s.service.Plans(r.Context())
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) createPlan(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var req model.CreatePlanRequest
	if err := decode(r, &req); err != nil {
		fail(w, err)
		return
	}
	value, err := s.service.CreatePlan(r.Context(), req)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusCreated, value)
}
func (s *Server) plan(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	id := r.URL.Query().Get("id")
	value, err := s.service.AssessPlan(r.Context(), id)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) validatePlan(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	id := r.URL.Query().Get("id")
	value, err := s.service.ValidatePlan(r.Context(), id)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) activatePlan(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	id := r.URL.Query().Get("id")
	value, err := s.service.ActivatePlan(r.Context(), id)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) evaluateEvent(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	var req model.EvaluateEventRequest
	if err := decode(r, &req); err != nil {
		fail(w, err)
		return
	}
	value, err := s.service.EvaluateEvent(r.Context(), req)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusCreated, value)
}
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	value, err := s.service.Events(r.Context(), r.URL.Query().Get("feeder_id"), 100)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) eventSummary(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	value, err := s.service.EventSummary(r.Context(), r.URL.Query().Get("feeder_id"))
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	value, err := s.service.Stats(r.Context())
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) recover(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	if err := s.service.Recover(r.Context()); err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, map[string]string{"status": "recovered"})
}
func (s *Server) active(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	id := r.URL.Query().Get("feeder_id")
	value, err := s.service.Store().Active(r.Context(), id)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) capacity(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	value, err := s.service.Stats(r.Context())
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, map[string]int{"feeders": value.Feeders, "zones": value.Zones, "relays": value.Relays, "plans": value.Plans})
}
func (s *Server) audit(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	id := r.URL.Query().Get("plan_id")
	value, err := s.service.Plan(r.Context(), id)
	if err != nil {
		fail(w, err)
		return
	}
	write(w, http.StatusOK, value)
}
func (s *Server) version(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	write(w, http.StatusOK, map[string]string{"service": "gridguard", "go": "1.26.3"})
}
func (s *Server) help(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	write(w, http.StatusOK, map[string]any{"service": "配电保护整定协调服务", "routes": []string{"/v1/feeders/create", "/v1/zones/create", "/v1/relays/create", "/v1/plans/create", "/v1/plans/validate", "/v1/plans/activate", "/v1/events/evaluate"}})
}
