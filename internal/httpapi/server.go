package httpapi

import (
	"github.com/rr173/task150-gridguard/internal/service"
	"net/http"
)

type Server struct {
	service *service.Service
	mux     *http.ServeMux
}

func New(s *service.Service) *Server {
	server := &Server{service: s, mux: http.NewServeMux()}
	server.routes()
	return server
}
func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.health)
	s.mux.HandleFunc("/readyz", s.ready)
	s.mux.HandleFunc("/v1/feeders", s.feeders)
	s.mux.HandleFunc("/v1/feeders/create", s.createFeeder)
	s.mux.HandleFunc("/v1/zones", s.zones)
	s.mux.HandleFunc("/v1/zones/create", s.createZone)
	s.mux.HandleFunc("/v1/relays", s.relays)
	s.mux.HandleFunc("/v1/relays/create", s.createRelay)
	s.mux.HandleFunc("/v1/plans", s.plans)
	s.mux.HandleFunc("/v1/plans/create", s.createPlan)
	s.mux.HandleFunc("/v1/plans/get", s.plan)
	s.mux.HandleFunc("/v1/plans/validate", s.validatePlan)
	s.mux.HandleFunc("/v1/plans/activate", s.activatePlan)
	s.mux.HandleFunc("/v1/events/evaluate", s.evaluateEvent)
	s.mux.HandleFunc("/v1/events", s.events)
	s.mux.HandleFunc("/v1/events/summary", s.eventSummary)
	s.mux.HandleFunc("/v1/stats", s.stats)
	s.mux.HandleFunc("/v1/recover", s.recover)
	s.mux.HandleFunc("/v1/active", s.active)
	s.mux.HandleFunc("/v1/capacity", s.capacity)
	s.mux.HandleFunc("/v1/audit", s.audit)
	s.mux.HandleFunc("/v1/version", s.version)
	s.mux.HandleFunc("/v1/help", s.help)
}
