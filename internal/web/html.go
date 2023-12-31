package web

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/rs/xid"
	"go.albinodrought.com/creamy-board/internal/cfg"
	"go.albinodrought.com/creamy-board/internal/log"
	"go.albinodrought.com/creamy-board/internal/repo"
	"go.albinodrought.com/creamy-board/internal/thumbnailer"
	"go.albinodrought.com/creamy-board/internal/web/tmpl"
)

func htmlUnhandled(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

var (
	UIErrorThreadBodyRequired     = "thread_body_required"
	UIErrorThreadFileRequired     = "thread_file_required"
	UIErrorUnsupportedMimetype    = "unsupported_mimetype"
	UIErrorPostBodyOrFileRequired = "post_body_or_file_required"
)

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

	if len(boards) == 1 {
		http.Redirect(w, r, fmt.Sprintf("/%v/index.html", boards[0].Slug), http.StatusFound)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	tmpl.ListBoards(boards).Render(r.Context(), w)
}

func (wp *HTMLWebPortal) ListBoardThreads(w http.ResponseWriter, r *http.Request, boardSlug string, page int) {
	boardRecentThreads, err := wp.Repo.ShowBoardListRecentThreads(r.Context(), boardSlug, page)
	if err != nil {
		log.Warnf("failed to list board %v threads: %v", boardSlug, err)
		htmlUnhandled(w) // todo: could be 404 (boardslug)
		return
	}

	var errorText string
	prevError := r.URL.Query().Get("error")
	switch prevError {
	case UIErrorThreadBodyRequired:
		errorText = "Threads must have a body"
	case UIErrorThreadFileRequired:
		errorText = "Threads must have at least one file"
	case UIErrorUnsupportedMimetype:
		errorText = "Thread contained an unsupported filetype"
	}

	w.Header().Add("Content-Type", "text/html")
	tmpl.ShowBoardAndRecents(boardRecentThreads, page, errorText).Render(r.Context(), w)
}

func (wp *HTMLWebPortal) ShowThread(w http.ResponseWriter, r *http.Request, boardSlug string, threadID int) {
	boardFullThread, err := wp.Repo.ShowThread(r.Context(), boardSlug, threadID)
	if err != nil {
		log.Warnf("failed to show board %v thread %v: %v", boardSlug, threadID, err)
		htmlUnhandled(w) // todo: could be 404 (bnoardSlug, threadID)
		return
	}

	var errorText string
	prevError := r.URL.Query().Get("error")
	switch prevError {
	case UIErrorPostBodyOrFileRequired:
		errorText = "Posts must have a body or file"
	case UIErrorUnsupportedMimetype:
		errorText = "Post contained an unsupported filetype"
	}

	w.Header().Add("Content-Type", "text/html")
	tmpl.ShowFullThread(boardFullThread, errorText).Render(r.Context(), w)
}

var acceptedMimes = []string{
	"image/jpeg",
	"image/gif",
	"image/png",
}

func htmlMimeTypeAllowed(mime *mimetype.MIME) bool {
	mimeStr := mime.String()
	for _, acceptedMime := range acceptedMimes {
		if acceptedMime == mimeStr {
			return true
		}
	}
	return false
}

var ErrMimeTypeNotAllowed = errors.New("mimetype not allowed")

func (wp *HTMLWebPortal) saveFormFile(r *http.Request, boardSlug string, key string) (*repo.SubmitPostFile, error) {
	formFile, formFileHeader, err := r.FormFile(key)
	if err != nil {
		if err == http.ErrMissingFile {
			// not submitted, normal
			return nil, nil
		}
		return nil, err
	}
	defer formFile.Close()

	mime, _ := mimetype.DetectReader(formFile)
	if mime == nil || !htmlMimeTypeAllowed(mime) {
		return nil, ErrMimeTypeNotAllowed
	}
	extension := strings.TrimPrefix(mime.Extension(), ".")
	if extension == "" {
		extension = "unknown"
	}

	filePath := path.Join("uf", path.Clean(boardSlug), xid.New().String())

	_, err = formFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	err = cfg.Storage.Write(r.Context(), filePath, formFile)
	if err != nil {
		return nil, err
	}

	postFile := repo.SubmitPostFile{
		Extension:    extension,
		MimeType:     mime.String(),
		Bytes:        int(formFileHeader.Size),
		OriginalName: formFileHeader.Filename,
		InternalPath: filePath,
	}

	_, err = formFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Warnf("error while seeking before thumbnail gen, skipping thumbnail: %v", err)
	} else {
		thumb, err := thumbnailer.Generate(formFile, mime)
		if err != nil {
			log.Warnf("error during thumbnail gen for %v, skipping: %v", mime.String(), err)
		} else {
			defer thumb.Close()
			thumbFilePath := filePath + ".thumb"
			err = cfg.Storage.Write(r.Context(), thumbFilePath, thumb)
			if err != nil {
				cfg.Storage.Delete(r.Context(), postFile.InternalPath)
				return nil, err
			}

			thumbExtension := strings.TrimPrefix(thumb.MIME.Extension(), ".")
			if thumbExtension == "" {
				thumbExtension = "unknown"
			}

			postFile.ThumbExtension = thumbExtension
			postFile.ThumbMimeType = thumb.MIME.String()
			postFile.ThumbBytes = thumb.Bytes
			postFile.ThumbInternalPath = thumbFilePath
		}
	}

	return &postFile, nil
}

func (wp *HTMLWebPortal) cleanupSubmit(req repo.SubmitPost) {
	for _, file := range req.Files {
		cfg.Storage.Delete(context.TODO(), file.InternalPath)
		if file.ThumbInternalPath != "" {
			cfg.Storage.Delete(context.TODO(), file.ThumbInternalPath)
		}
	}
}

func (wp *HTMLWebPortal) SubmitThread(w http.ResponseWriter, r *http.Request, boardSlug string) {
	req := repo.SubmitPost{
		Subject: r.FormValue("subject"),
		Author:  r.FormValue("author"),
		Body:    r.FormValue("body"),
		Files:   []repo.SubmitPostFile{},
	}

	if req.Author == "" {
		req.Author = "Anonymous"
	}

	if req.Body == "" {
		http.Redirect(w, r, fmt.Sprintf("/%v/index.html?error=%v", boardSlug, UIErrorThreadBodyRequired), http.StatusFound)
		return
	}

	created := false
	// clean up temp files if req fails
	defer func() {
		if !created {
			wp.cleanupSubmit(req)
		}
	}()

	for _, fileKey := range []string{"file1", "file2", "file3"} {
		postFile, err := wp.saveFormFile(r, boardSlug, fileKey)
		if err != nil {
			if err == ErrMimeTypeNotAllowed {
				http.Redirect(w, r, fmt.Sprintf("/%v/index.html?error=%v", boardSlug, UIErrorUnsupportedMimetype), http.StatusFound)
				return
			}
			log.Warnf("failed to save form file %v: %v", fileKey, err)
			htmlUnhandled(w)
			return
		}
		if postFile != nil {
			req.Files = append(req.Files, *postFile)
		}
	}

	if len(req.Files) < 1 {
		http.Redirect(w, r, fmt.Sprintf("/%v/index.html?error=%v", boardSlug, UIErrorThreadFileRequired), http.StatusFound)
		return
	}

	threadID, err := wp.Repo.SubmitThread(r.Context(), boardSlug, req)
	if err != nil {
		log.Warnf("failed to create board %v thread %+v: %v", boardSlug, req, err)
		htmlUnhandled(w) // todo: could be 404 (boardSlug)
		return
	}

	created = true
	http.Redirect(w, r, fmt.Sprintf("/%v/res/%v.html", boardSlug, threadID), http.StatusFound)
}

func (wp *HTMLWebPortal) SubmitThreadPost(w http.ResponseWriter, r *http.Request, boardSlug string, threadID int) {
	req := repo.SubmitPost{
		Subject: r.FormValue("subject"),
		Author:  r.FormValue("author"),
		Body:    r.FormValue("body"),
		Files:   []repo.SubmitPostFile{},
	}

	if req.Author == "" {
		req.Author = "Anonymous"
	}

	created := false
	// clean up temp files if req fails
	defer func() {
		if !created {
			wp.cleanupSubmit(req)
		}
	}()

	for _, fileKey := range []string{"file1", "file2", "file3"} {
		postFile, err := wp.saveFormFile(r, boardSlug, fileKey)
		if err != nil {
			if err == ErrMimeTypeNotAllowed {
				http.Redirect(w, r, fmt.Sprintf("/%v/res/%v.html?error=%v", boardSlug, threadID, UIErrorUnsupportedMimetype), http.StatusFound)
				return
			}
			log.Warnf("failed to save form file %v: %v", fileKey, err)
			htmlUnhandled(w)
			return
		}
		if postFile != nil {
			req.Files = append(req.Files, *postFile)
		}
	}

	if req.Body == "" && len(req.Files) < 1 {
		http.Redirect(w, r, fmt.Sprintf("/%v/res/%v.html?error=%v", boardSlug, threadID, UIErrorPostBodyOrFileRequired), http.StatusFound)
		return
	}

	postID, err := wp.Repo.SubmitThreadPost(r.Context(), boardSlug, threadID, req)
	if err != nil {
		log.Warnf("failed to create board %v thread %v post %+v: %v", boardSlug, threadID, req, err)
		htmlUnhandled(w) // todo: could be 404 (boardSlug, threadID)
		return
	}

	created = true
	http.Redirect(w, r, fmt.Sprintf("/%v/res/%v.html#%v", boardSlug, threadID, postID), http.StatusFound)
}

var _ WebPortal = &HTMLWebPortal{}
