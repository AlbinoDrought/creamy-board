package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.albinodrought.com/creamy-board/internal/cfg"
	"go.albinodrought.com/creamy-board/internal/repo"
)

func Router() http.Handler {
	r := chi.NewRouter()

	repo := repo.DBRepo{
		Querier: cfg.Querier,
	}

	htmlPortal := HTMLWebPortal{
		Repo: &repo,
	}

	r.Get("/", htmlPortal.ListBoards)
	r.Get("/index.html", htmlPortal.ListBoards)
	r.Get("/{boardSlug}/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("/%v/index.html", chi.URLParam(r, "boardSlug")), http.StatusFound)
	})
	r.Get("/{boardSlug}/index.html", func(w http.ResponseWriter, r *http.Request) {
		htmlPortal.ListBoardThreads(w, r, chi.URLParam(r, "boardSlug"), 1)
	})
	r.Get("/{boardSlug}/{page}.html", func(w http.ResponseWriter, r *http.Request) {
		pageStr := chi.URLParam(r, "page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			page = 0
		}
		if page == 1 {
			http.Redirect(w, r, fmt.Sprintf("/%v/index.html", chi.URLParam(r, "boardSlug")), http.StatusFound)
			return
		}

		htmlPortal.ListBoardThreads(w, r, chi.URLParam(r, "boardSlug"), page)
	})
	r.Get("/{boardSlug}/res/{threadID}.html", func(w http.ResponseWriter, r *http.Request) {
		threadIDStr := chi.URLParam(r, "threadID")
		threadID, err := strconv.Atoi(threadIDStr)
		if err != nil || threadID < 0 {
			threadID = 0 // let handler show 404
		}

		htmlPortal.ShowThread(w, r, chi.URLParam(r, "boardSlug"), threadID)
	})

	jsonPortal := JSONWebPortal{
		Repo: &repo,
	}

	r.Get("/index.json", jsonPortal.ListBoards)
	r.Get("/{boardSlug}/index.json", func(w http.ResponseWriter, r *http.Request) {
		jsonPortal.ListBoardThreads(w, r, chi.URLParam(r, "boardSlug"), 1)
	})
	r.Get("/{boardSlug}/{page}.json", func(w http.ResponseWriter, r *http.Request) {
		pageStr := chi.URLParam(r, "page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			page = 0
		}
		if page == 1 {
			http.Redirect(w, r, fmt.Sprintf("/%v/index.json", chi.URLParam(r, "boardSlug")), http.StatusFound)
			return
		}

		jsonPortal.ListBoardThreads(w, r, chi.URLParam(r, "boardSlug"), page)
	})
	r.Get("/{boardSlug}/res/{threadID}.json", func(w http.ResponseWriter, r *http.Request) {
		threadIDStr := chi.URLParam(r, "threadID")
		threadID, err := strconv.Atoi(threadIDStr)
		if err != nil || threadID < 0 {
			threadID = 0 // let handler show 404
		}
		jsonPortal.ShowThread(w, r, chi.URLParam(r, "boardSlug"), threadID)
	})

	return r
}
