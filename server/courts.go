package server

import (
	"fmt"
	"net/http"

	"github.com/broothie/dillbook/model"
	"github.com/go-chi/chi/v5"
)

func (s *Server) NewCourt(w http.ResponseWriter, r *http.Request) {
	var court model.Court
	if err := s.decodeForm(r, &court); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.render.HTML(w, http.StatusOK, "courts/new", map[string]any{"court": court})
}

func (s *Server) CreateCourt(w http.ResponseWriter, r *http.Request) {
	var court model.Court
	if err := s.decodeForm(r, &court); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.DB.Create(court).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/courts/%v", court.ID), http.StatusSeeOther)
}

func (s *Server) ShowCourt(w http.ResponseWriter, r *http.Request) {
	var court model.Court
	if err := s.DB.Find(&court, chi.URLParam(r, "courtID")).Error; err != nil {
		http.Redirect(w, r, "/courts", http.StatusPermanentRedirect)
		return
	}

	s.render.HTML(w, http.StatusOK, "courts/show", map[string]any{"court": court})
}
