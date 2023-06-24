package web

import (
	"net/http"

	"go.albinodrought.com/creamy-board/internal/log"
	"go.albinodrought.com/creamy-board/internal/repo"
	"go.albinodrought.com/creamy-board/internal/web/tmpl"
)

func htmlUnhandled(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

type HTMLWebPortal struct {
	Repo repo.CreamyBoard
}

func (wp *HTMLWebPortal) ListBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := wp.Repo.ListBoards(r.Context())
	if err != nil {
		log.Warnf("failed to list boards: %v", err)
		htmlUnhandled(w)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	tmpl.ListBoards(boards).Render(r.Context(), w)
}

func (wp *HTMLWebPortal) ListBoardThreads(w http.ResponseWriter, r *http.Request, boardSlug string, page int) {
	boardRecentThreads, err := wp.Repo.ShowBoardListRecenthreads(r.Context(), boardSlug, page)
	if err != nil {
		log.Warnf("failed to list board %v threads: %v", boardSlug, err)
		htmlUnhandled(w) // todo: could be 404 (boardslug)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	tmpl.ShowBoardAndRecents(boardRecentThreads).Render(r.Context(), w)
}

func (wp *HTMLWebPortal) ShowThread(w http.ResponseWriter, r *http.Request, boardSlug string, threadID int) {
	boardFullThread, err := wp.Repo.ShowThread(r.Context(), boardSlug, threadID)
	if err != nil {
		log.Warnf("failed to show board %v thread %v: %v", boardSlug, threadID, err)
		htmlUnhandled(w) // todo: could be 404 (bnoardSlug, threadID)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	tmpl.ShowFullThread(boardFullThread).Render(r.Context(), w)
}

var _ WebPortal = &HTMLWebPortal{}
