package web

import "net/http"

type WebPortal interface {
	ListBoards(w http.ResponseWriter, r *http.Request)
	ShowBoard(w http.ResponseWriter, r *http.Request, boardSlug string)
	ListBoardThreads(w http.ResponseWriter, r *http.Request, boardSlug string, page int)
	ShowThread(w http.ResponseWriter, r *http.Request, boardSlug string, threadID int)
}
