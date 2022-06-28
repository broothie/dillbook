package application

import (
	"fmt"
	"net/http"

	"github.com/broothie/dillbook/model"
	"github.com/go-chi/chi/v5"
)

func (a *Application) Index(w http.ResponseWriter, r *http.Request) {
	a.render.HTML(w, http.StatusOK, "index", nil)
}

func (a *Application) NewCourt(w http.ResponseWriter, r *http.Request) {
	a.render.HTML(w, http.StatusOK, "courts/new", model.Court{Name: r.FormValue("name")})
}

func (a *Application) CreateCourt(w http.ResponseWriter, r *http.Request) {
	court := &model.Court{Name: r.FormValue("name")}
	if err := a.DB.Create(court).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/courts/%v", court.ID), http.StatusSeeOther)
}

func (a *Application) ShowCourt(w http.ResponseWriter, r *http.Request) {
	var court model.Court
	if err := a.DB.Find(&court, chi.URLParam(r, "courtID")).Error; err != nil {
		http.Redirect(w, r, "/courts", http.StatusPermanentRedirect)
		return
	}

	a.render.HTML(w, http.StatusOK, "courts/show", court)
}
