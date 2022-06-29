package server

import (
	"fmt"
	"net/http"

	"github.com/broothie/dillbook/model"
	"github.com/go-chi/chi/v5"
)

func (s *Server) newCourt(w http.ResponseWriter, r *http.Request) {
	var court model.Court
	if err := s.decodeForm(r, &court); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	s.render.HTML(w, http.StatusOK, "courts/new", map[string]any{
		"court":      court,
		"locationID": chi.URLParam(r, "locationID"),
	})
}

func (s *Server) createCourt(w http.ResponseWriter, r *http.Request) {
	court := model.Court{LocationID: chi.URLParam(r, "locationID")}
	if err := s.decodeForm(r, &court); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := s.validate.Struct(court); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := s.DB.WithContext(r.Context()).Create(&court).Error; err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/courts/%v", court.ID), http.StatusSeeOther)
}

func (s *Server) showCourt(w http.ResponseWriter, r *http.Request) {
	var court model.Court
	if err := s.DB.WithContext(r.Context()).Preload("Location").Find(&court, "id = ?", chi.URLParam(r, "courtID")).Error; err != nil {
		http.Redirect(w, r, "/courts", http.StatusPermanentRedirect)
		return
	}

	s.render.HTML(w, http.StatusOK, "courts/show", map[string]any{
		"locationID": chi.URLParam(r, "locationID"),
		"court":      court,
	})
}
