package web

import (
	"fmt"
	"io"
	"net/http"

	"go.albinodrought.com/creamy-board/internal/cfg"
	"go.albinodrought.com/creamy-board/internal/db/queries"
	"go.albinodrought.com/creamy-board/internal/log"
)

func fileUnhandled(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func ServeFile(w http.ResponseWriter, r *http.Request, boardSlug string, threadID int, postID int, index int, extension string, thumb bool) {
	dbBoard, err := cfg.Querier.ShowBoardFromSlug(r.Context(), boardSlug)
	if err != nil {
		log.Warnf("failed to show board %v: %v", boardSlug, err)
		fileUnhandled(w) // todo: 404
		return
	}

	var (
		path  string
		ext   string
		mime  string
		bytes int
	)

	if thumb {
		dbFile, err := cfg.Querier.ShowFileThumb(r.Context(), queries.ShowFileThumbParams{
			BoardID:  dbBoard.BoardID,
			ThreadID: threadID,
			PostID:   postID,
			Idx:      int16(index),
		})
		if err != nil {
			log.Warnf("failed to show file thumb %v/%v-%v-%v.%v: %v", boardSlug, threadID, postID, index, extension, err)
			fileUnhandled(w) // todo: 404
			return
		}
		path = dbFile.ThumbPath.String
		ext = dbFile.ThumbExtension.String
		mime = dbFile.ThumbMimetype.String
		bytes = int(*dbFile.ThumbBytes)
	} else {
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
		path = dbFile.Path.String
		ext = dbFile.Extension.String
		mime = dbFile.Mimetype.String
		bytes = int(dbFile.Bytes)
	}

	if ext != extension {
		log.Warnf("failed to show file %v/%v-%v-%v.%v: actually has extension %v", boardSlug, threadID, postID, index, extension, ext)
		fileUnhandled(w) // todo: 404
		return
	}

	handle, err := cfg.Storage.Read(r.Context(), path)
	if err != nil {
		log.Warnf("failed to read file %v/%v-%v-%v.%v: %v", boardSlug, threadID, postID, index, extension, err)
		fileUnhandled(w)
		return
	}
	defer handle.Close()

	w.Header().Add("Content-Length", fmt.Sprintf("%v", bytes))
	w.Header().Add("Content-Type", mime)
	w.Header().Add("Content-Disposition", "inline")
	w.Header().Add("Cache-Control", "public, max-age=31536000, immutable")
	io.Copy(w, handle)
}
