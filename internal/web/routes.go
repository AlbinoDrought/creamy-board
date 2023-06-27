package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi"
	"go.albinodrought.com/creamy-board/internal/cfg"
	"go.albinodrought.com/creamy-board/internal/repo"
	"go.albinodrought.com/creamy-board/internal/web/static"
)

func Router() http.Handler {
	r := chi.NewRouter()

	fileServer := http.FileServer(http.FS(static.FS))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "public, max-age=86400, stale-while-revalidate")
		fileServer.ServeHTTP(w, r)
	})

	repo := repo.DBRepo{
		Querier: cfg.Querier,
	}

	r.Get("/{boardSlug}/src/{threadID}-{postID}-{index}.{extension}", func(w http.ResponseWriter, r *http.Request) {
		threadIDStr := chi.URLParam(r, "threadID")
		threadID, err := strconv.Atoi(threadIDStr)
		if err != nil || threadID < 0 {
			threadID = 0 // let handler show 404
		}

		postIDStr := chi.URLParam(r, "postID")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil || postID < 0 {
			postID = 0 // let handler show 404
		}

		indexStr := chi.URLParam(r, "index")
		index, err := strconv.Atoi(indexStr)
		if err != nil || postID < 0 {
			index = 0 // let handler show 404
		}

		ServeFile(w, r, chi.URLParam(r, "boardSlug"), threadID, postID, index, chi.URLParam(r, "extension"), false)
	})

	r.Get("/{boardSlug}/thumb/{threadID}-{postID}-{index}.{extension}", func(w http.ResponseWriter, r *http.Request) {
		threadIDStr := chi.URLParam(r, "threadID")
		threadID, err := strconv.Atoi(threadIDStr)
		if err != nil || threadID < 0 {
			threadID = 0 // let handler show 404
		}

		postIDStr := chi.URLParam(r, "postID")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil || postID < 0 {
			postID = 0 // let handler show 404
		}

		indexStr := chi.URLParam(r, "index")
		index, err := strconv.Atoi(indexStr)
		if err != nil || postID < 0 {
			index = 0 // let handler show 404
		}

		ServeFile(w, r, chi.URLParam(r, "boardSlug"), threadID, postID, index, chi.URLParam(r, "extension"), true)
	})

	htmlPortal := HTMLWebPortal{
		Repo: &repo,
	}

	r.Get("/", htmlPortal.ListBoards)
	r.Get("/index.html", htmlPortal.ListBoards)
	// todo: redirect /{boardSlug}, but only if it is a valid board (otherwise ruins /favicon.ico routing)
	r.Get("/{boardSlug}/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("/%v/index.html", chi.URLParam(r, "boardSlug")), http.StatusFound)
	})
	r.Get("/{boardSlug}/index.html", func(w http.ResponseWriter, r *http.Request) {
		htmlPortal.ListBoardThreads(w, r, chi.URLParam(r, "boardSlug"), 1)
	})
	r.Post("/{boardSlug}/index.html", func(w http.ResponseWriter, r *http.Request) {
		htmlPortal.SubmitThread(w, r, chi.URLParam(r, "boardSlug"))
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
	r.Post("/{boardSlug}/res/{threadID}.html", func(w http.ResponseWriter, r *http.Request) {
		threadIDStr := chi.URLParam(r, "threadID")
		threadID, err := strconv.Atoi(threadIDStr)
		if err != nil || threadID < 0 {
			threadID = 0 // let handler show 404
		}

		htmlPortal.SubmitThreadPost(w, r, chi.URLParam(r, "boardSlug"), threadID)
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

	return gziphandler.GzipHandler(r)
}
