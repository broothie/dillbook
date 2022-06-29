package server

import (
	"fmt"
	"net/http"

	"github.com/broothie/dillbook/model"
	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
)

func (s *Server) indexLocations(w http.ResponseWriter, r *http.Request) {
	const gridWidth = 3

	var locations []model.Location
	if err := s.DB.WithContext(r.Context()).Preload("Courts").Find(&locations).Error; err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	s.render.HTML(w, http.StatusOK, "locations/index", map[string]any{
		"gridWidth":        gridWidth,
		"locations":        locations,
		"chunkedLocations": lo.Chunk(locations, gridWidth),
	})
}

func (s *Server) newLocation(w http.ResponseWriter, r *http.Request) {
	var location model.Location
	if err := s.decodeForm(r, &location); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	s.render.HTML(w, http.StatusOK, "locations/new", map[string]any{"location": location})
}

func (s *Server) createLocation(w http.ResponseWriter, r *http.Request) {
	var location model.Location
	if err := s.decodeForm(r, &location); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := s.validate.Struct(location); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := s.DB.WithContext(r.Context()).Create(&location).Error; err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/locations/%v", location.ID), http.StatusSeeOther)
}

func (s *Server) showLocation(w http.ResponseWriter, r *http.Request) {
	var location model.Location
	if err := s.DB.WithContext(r.Context()).Preload("Courts").Find(&location, "id = ?", chi.URLParam(r, "locationID")).Error; err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	s.render.HTML(w, http.StatusOK, "locations/show", map[string]any{"location": location})
}
