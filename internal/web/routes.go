package web

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Router() http.Handler {
	r := chi.NewRouter()

	jsonPortal := JSONWebPortal{}

	r.Get("/index.json", jsonPortal.ListBoards)
	r.Get("/{boardSlug}/index.json", func(w http.ResponseWriter, r *http.Request) {
		jsonPortal.ShowBoard(w, r, chi.URLParam(r, "boardSlug"))
	})

	return r
}
