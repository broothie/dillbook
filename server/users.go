package server

import (
	"fmt"
	"net/http"

	"github.com/broothie/dillbook/model"
	"github.com/go-chi/chi/v5"
)

func (s *Server) NewUser(w http.ResponseWriter, r *http.Request) {
	s.render.HTML(w, http.StatusOK, "users/new", model.User{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	})
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}

	if err := s.DB.Create(user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%v", user.ID), http.StatusSeeOther)
}

func (s *Server) ShowUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := s.DB.Find(&user, chi.URLParam(r, "userID")).Error; err != nil {
		http.Redirect(w, r, "/users", http.StatusPermanentRedirect)
		return
	}

	s.render.HTML(w, http.StatusOK, "users/show", user)
}
