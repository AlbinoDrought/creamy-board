package web

import (
	"encoding/json"
	"net/http"
	"time"

	"go.albinodrought.com/creamy-board/internal/cfg"
	"go.albinodrought.com/creamy-board/internal/db/queries"
	"go.albinodrought.com/creamy-board/internal/log"
)

type Board struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Tagline string `json:"tagline"`
}

func (b *Board) FromList(row queries.ListBoardsRow) {
	b.Slug = row.Slug.String
	b.Title = row.Title.String
	b.Tagline = row.Tagline.String
}

func (b *Board) FromShow(row queries.ShowBoardFromSlugRow) {
	b.Slug = row.Slug.String
	b.Title = row.Title.String
	b.Tagline = row.Tagline.String
}

type Thread struct {
	ID        uint64 `json:"id"`
	CreatedAt string `json:"created_at"`
	BumpedAt  string `json:"bumped_at"`
}

type Post struct {
	ID        uint64 `json:"id"`
	CreatedAt string `json:"created_at"`
	Author    string `json:"author"`
	Body      string `json:"body"`
}

type RecentThread struct {
	Thread      Thread `json:"thread"`
	MainPost    Post   `json:"main_post"`
	RecentPosts []Post `json:"recent_posts"`
}

func (rt *RecentThread) FromActive(row queries.ListActiveBoardThreadsRow) {
	rt.Thread.ID = uint64(*row.ThreadID)
	rt.Thread.CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
	rt.Thread.BumpedAt = row.BumpedAt.Time.Format(time.RFC3339)

	rt.MainPost.ID = rt.Thread.ID
	rt.MainPost.CreatedAt = rt.Thread.CreatedAt
	rt.MainPost.Author = row.Author.String
	rt.MainPost.Body = *row.Body
}

func (rt *RecentThread) LoadRecents(rows []queries.ListThreadRecentPostsRow) {
	rt.RecentPosts = make([]Post, len(rows))
	for i, row := range rows {
		rt.RecentPosts[i].ID = uint64(*row.PostID)
		rt.RecentPosts[i].CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
		rt.RecentPosts[i].Author = row.Author.String
		rt.RecentPosts[i].Body = *row.Body
	}
}

type FullThread struct {
	Thread   Thread `json:"thread"`
	MainPost Post   `json:"main_post"`
	AllPosts []Post `json:"all_posts"`
}

func (ft *FullThread) FromShow(row queries.ShowThreadRow) {
	ft.Thread.ID = uint64(*row.ThreadID)
	ft.Thread.CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
	ft.Thread.BumpedAt = row.BumpedAt.Time.Format(time.RFC3339)

	ft.MainPost.ID = ft.Thread.ID
	ft.MainPost.CreatedAt = ft.Thread.CreatedAt
	ft.MainPost.Author = row.Author.String
	ft.MainPost.Body = *row.Body
}

func (ft *FullThread) LoadAll(rows []queries.ListThreadPostsRow) {
	ft.AllPosts = make([]Post, len(rows))
	for i, row := range rows {
		ft.AllPosts[i].ID = uint64(*row.PostID)
		ft.AllPosts[i].CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
		ft.AllPosts[i].Author = row.Author.String
		ft.AllPosts[i].Body = *row.Body
	}
}

func unhandled(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

type JSONWebPortal struct{}

func (wp *JSONWebPortal) ListBoards(w http.ResponseWriter, r *http.Request) {
	dbBoards, err := cfg.Querier.ListBoards(r.Context())
	if err != nil {
		log.Warnf("failed to list boards: %v", err)
		unhandled(w)
		return
	}

	boards := make([]Board, len(dbBoards))
	for i := range boards {
		boards[i].FromList(dbBoards[i])
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&boards)
}

func (wp *JSONWebPortal) ShowBoard(w http.ResponseWriter, r *http.Request, boardSlug string) {
	dbBoard, err := cfg.Querier.ShowBoardFromSlug(r.Context(), boardSlug)
	if err != nil {
		log.Warnf("failed to show board %v: %v", boardSlug, err)
		unhandled(w) // todo: could be 404
		return
	}

	board := Board{}
	board.FromShow(dbBoard)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&board)
}

func (wp *JSONWebPortal) ListBoardThreads(w http.ResponseWriter, r *http.Request, boardSlug string, page int) {
	dbBoard, err := cfg.Querier.ShowBoardFromSlug(r.Context(), boardSlug)
	if err != nil {
		log.Warnf("failed to show board %v: %v", boardSlug, err)
		unhandled(w) // todo: could be 404
		return
	}

	dbThreads, err := cfg.Querier.ListActiveBoardThreads(r.Context(), queries.ListActiveBoardThreadsParams{
		BoardID: dbBoard.BoardID,
		Limit:   15,
		Offset:  (page - 1) * 15,
	})
	if err != nil {
		log.Warnf("failed to list active threads for board %v: %v", boardSlug, err)
		unhandled(w)
		return
	}

	threadIDs := make([]int, len(dbThreads))
	for i := range dbThreads {
		threadIDs[i] = *dbThreads[i].ThreadID
	}

	dbRecentPosts, err := cfg.Querier.ListThreadRecentPosts(r.Context(), dbBoard.BoardID, threadIDs)
	if err != nil {
		log.Warnf("failed to list recent thread posts for board %v: %v", boardSlug, err)
		unhandled(w)
		return
	}

	dbRecentPostsByThreads := make(map[int][]queries.ListThreadRecentPostsRow, len(dbThreads))
	for _, dbRecentPost := range dbRecentPosts {
		dbRecentPostsByThread, ok := dbRecentPostsByThreads[*dbRecentPost.ThreadID]
		if !ok {
			dbRecentPostsByThread = []queries.ListThreadRecentPostsRow{}
		}
		dbRecentPostsByThread = append(dbRecentPostsByThread, dbRecentPost)
		dbRecentPostsByThreads[*dbRecentPost.ThreadID] = dbRecentPostsByThread
	}

	recentThreads := make([]RecentThread, len(dbThreads))
	for i := range dbThreads {
		recentThreads[i].FromActive(dbThreads[i])
		dbRecentPosts, ok := dbRecentPostsByThreads[*dbThreads[i].ThreadID]
		if ok {
			recentThreads[i].LoadRecents(dbRecentPosts)
		} else {
			recentThreads[i].LoadRecents([]queries.ListThreadRecentPostsRow{})
		}
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&recentThreads)
}

func (wp *JSONWebPortal) ShowThread(w http.ResponseWriter, r *http.Request, boardSlug string, threadID int) {
	dbBoard, err := cfg.Querier.ShowBoardFromSlug(r.Context(), boardSlug)
	if err != nil {
		log.Warnf("failed to show board %v: %v", boardSlug, err)
		unhandled(w) // todo: could be 404
		return
	}

	dbThread, err := cfg.Querier.ShowThread(r.Context(), dbBoard.BoardID, threadID)
	if err != nil {
		log.Warnf("failed to show board %v thread %v: %v", boardSlug, threadID, err)
		unhandled(w) // todo: could be 404
		return
	}

	dbPosts, err := cfg.Querier.ListThreadPosts(r.Context(), dbBoard.BoardID, threadID)
	if err != nil {
		log.Warnf("failed to list board %v thread %v posts: %v", dbBoard.BoardID, threadID, err)
		unhandled(w)
		return
	}

	fullThread := FullThread{}
	fullThread.FromShow(dbThread)
	fullThread.LoadAll(dbPosts)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&fullThread)
}

var _ WebPortal = &JSONWebPortal{}
