package web

import (
	"io"
	"net/http"

	"go.albinodrought.com/creamy-board/internal/cfg"
	"go.albinodrought.com/creamy-board/internal/db/queries"
	"go.albinodrought.com/creamy-board/internal/log"
)

func fileUnhandled(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func ServeFile(w http.ResponseWriter, r *http.Request, boardSlug string, threadID int, postID int, index int, extension string) {
	dbBoard, err := cfg.Querier.ShowBoardFromSlug(r.Context(), boardSlug)
	if err != nil {
		log.Warnf("failed to show board %v: %v", boardSlug, err)
		fileUnhandled(w) // todo: 404
		return
	}

	dbFile, err := cfg.Querier.ShowFile(r.Context(), queries.ShowFileParams{
		BoardID:  dbBoard.BoardID,
		ThreadID: threadID,
		PostID:   postID,
		Idx:      int16(index),
	})
	if err != nil {
		log.Warnf("failed to show file %v/%v-%v-%v.%v: %v", boardSlug, threadID, postID, index, extension, err)
		fileUnhandled(w) // todo: 404
		return
	}

	if dbFile.Extension.String != extension {
		log.Warnf("failed to show file %v/%v-%v-%v.%v: actually has extension %v", boardSlug, threadID, postID, index, extension, dbFile.Extension.String)
		fileUnhandled(w) // todo: 404
		return
	}

	handle, err := cfg.Storage.Read(r.Context(), dbFile.Path.String)
	if err != nil {
		log.Warnf("failed to read file %v/%v-%v-%v.%v: %v", boardSlug, threadID, postID, index, extension, err)
		fileUnhandled(w)
		return
	}
	defer handle.Close()

	w.Header().Add("Content-Type", dbFile.Mimetype.String)
	w.Header().Add("Content-Disposition", "inline")
	w.Header().Add("Cache-Control", "public, max-age=31536000, immutable")
	io.Copy(w, handle)
}
