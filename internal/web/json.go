package web

import (
	"encoding/json"
	"net/http"

	"go.albinodrought.com/creamy-board/internal/cfg"
	"go.albinodrought.com/creamy-board/internal/log"
)

type JSONWebPortal struct{}

func (wp *JSONWebPortal) ListBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := cfg.Querier.ListBoards(r.Context())
	if err != nil {
		log.Warnf("failed to list boards: %v", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&boards)
}

func (wp *JSONWebPortal) ShowBoard(w http.ResponseWriter, r *http.Request, slug string) {
	board, err := cfg.Querier.ShowBoardFromSlug(r.Context(), slug)
	if err != nil {
		log.Warnf("failed to show board %v: %v", slug, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&board)
}

var _ WebPortal = &JSONWebPortal{}
