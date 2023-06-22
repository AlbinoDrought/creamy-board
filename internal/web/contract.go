package web

import "net/http"

type WebPortal interface {
	ListBoards(w http.ResponseWriter, r *http.Request)
	ShowBoard(w http.ResponseWriter, r *http.Request, slug string)
}
