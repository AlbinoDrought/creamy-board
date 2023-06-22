package web

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func Router() http.Handler {
	r := chi.NewRouter()

	jsonPortal := JSONWebPortal{}

	r.Get("/boards.json", jsonPortal.ListBoards)
	r.Get("/boards/{boardSlug}/info.json", func(w http.ResponseWriter, r *http.Request) {
		jsonPortal.ShowBoard(w, r, chi.URLParam(r, "boardSlug"))
	})
	r.Get("/boards/{boardSlug}/threads.json", func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}
		jsonPortal.ListBoardThreads(w, r, chi.URLParam(r, "boardSlug"), page)
	})
	r.Get("/boards/{boardSlug}/threads/{threadID}.json", func(w http.ResponseWriter, r *http.Request) {
		threadIDStr := chi.URLParam(r, "threadID")
		threadID, err := strconv.Atoi(threadIDStr)
		if err != nil || threadID <= 0 {
			threadID = 0 // let handler show 404
		}
		jsonPortal.ShowThread(w, r, chi.URLParam(r, "boardSlug"), threadID)
	})

	return r
}
