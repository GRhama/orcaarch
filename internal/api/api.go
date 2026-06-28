package api

import (
	"encoding/json"
	"net/http"

	"orcaarch/internal/report"
	"orcaarch/internal/service"
)

type Server struct {
	mux       *http.ServeMux
	inventory []report.ReconciliationRow
	ledger    []report.LedgerRow
	risk      []report.RiskRow
}

func NewServer(result service.ProcessingResult) *Server {
	s := &Server{
		mux:       http.NewServeMux(),
		inventory: report.ToReconciliationRows(result.Reconciliations),
		ledger:    report.ToLedgerRows(result.LedgerEntries),
		risk:      report.ToRiskRows(result.Quarantine),
	}
	s.mux.HandleFunc("GET /api/v1/reports/inventory", s.handleInventory)
	s.mux.HandleFunc("GET /api/v1/reports/ledger", s.handleLedger)
	s.mux.HandleFunc("GET /api/v1/reports/risk", s.handleRisk)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (s *Server) handleInventory(w http.ResponseWriter, r *http.Request) { writeJSON(w, s.inventory) }
func (s *Server) handleLedger(w http.ResponseWriter, r *http.Request)    { writeJSON(w, s.ledger) }
func (s *Server) handleRisk(w http.ResponseWriter, r *http.Request)      { writeJSON(w, s.risk) }
