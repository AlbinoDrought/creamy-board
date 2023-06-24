package web

import (
	"encoding/json"
	"net/http"

	"go.albinodrought.com/creamy-board/internal/log"
	"go.albinodrought.com/creamy-board/internal/repo"
)

func unhandled(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

type JSONWebPortal struct {
	Repo repo.CreamyBoard
}

func (wp *JSONWebPortal) ListBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := wp.Repo.ListBoards(r.Context())
	if err != nil {
		log.Warnf("failed to list boards: %v", err)
		unhandled(w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&boards)
}

func (wp *JSONWebPortal) ShowBoard(w http.ResponseWriter, r *http.Request, boardSlug string) {
	board, err := wp.Repo.ShowBoard(r.Context(), boardSlug)
	if err != nil {
		log.Warnf("failed to show board %v: %v", boardSlug, err)
		unhandled(w) // todo: could be 404 (boardslug)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&board)
}

func (wp *JSONWebPortal) ListBoardThreads(w http.ResponseWriter, r *http.Request, boardSlug string, page int) {
	recentThreads, err := wp.Repo.ListBoardThreads(r.Context(), boardSlug, page)
	if err != nil {
		log.Warnf("failed to list board %v threads: %v", boardSlug, err)
		unhandled(w) // todo: could be 404 (boardslug)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&recentThreads)
}

func (wp *JSONWebPortal) ShowThread(w http.ResponseWriter, r *http.Request, boardSlug string, threadID int) {
	fullThread, err := wp.Repo.ShowThread(r.Context(), boardSlug, threadID)
	if err != nil {
		log.Warnf("failed to show board %v thread %v: %v", boardSlug, threadID, err)
		unhandled(w) // todo: could be 404 (bnoardSlug, threadID)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&fullThread)
}

var _ WebPortal = &JSONWebPortal{}
