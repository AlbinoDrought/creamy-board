package repo

import (
	"context"
	"errors"
	"time"

	"go.albinodrought.com/creamy-board/internal/db/queries"
)

func boardFromList(b *Board, row queries.ListBoardsRow) {
	b.Slug = row.Slug.String
	b.Title = row.Title.String
	b.Tagline = row.Tagline.String
}

func boardFromShow(b *Board, row queries.ShowBoardFromSlugRow) {
	b.Slug = row.Slug.String
	b.Title = row.Title.String
	b.Tagline = row.Tagline.String
}

func recentThreadFromActive(rt *RecentThread, row queries.ListActiveBoardThreadsRow) {
	rt.Thread.ID = uint64(*row.ThreadID)
	rt.Thread.CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
	rt.Thread.BumpedAt = row.BumpedAt.Time.Format(time.RFC3339)
	rt.Thread.Subject = row.Subject.String

	rt.MainPost.ID = rt.Thread.ID
	rt.MainPost.CreatedAt = rt.Thread.CreatedAt
	rt.MainPost.Author = row.Author.String
	rt.MainPost.Body = *row.Body
}

func recentThreadLoadRecentPosts(rt *RecentThread, rows []queries.ListThreadRecentPostsRow) {
	rt.RecentPosts = make([]Post, len(rows))
	for i, row := range rows {
		rt.RecentPosts[i].ID = uint64(*row.PostID)
		rt.RecentPosts[i].CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
		rt.RecentPosts[i].Author = row.Author.String
		rt.RecentPosts[i].Body = *row.Body
	}
}

func fullThreadFromShow(ft *FullThread, row queries.ShowThreadRow) {
	ft.Thread.ID = uint64(*row.ThreadID)
	ft.Thread.CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
	ft.Thread.BumpedAt = row.BumpedAt.Time.Format(time.RFC3339)
	ft.Thread.Subject = row.Subject.String

	ft.MainPost.ID = ft.Thread.ID
	ft.MainPost.CreatedAt = ft.Thread.CreatedAt
	ft.MainPost.Author = row.Author.String
	ft.MainPost.Body = *row.Body
}

func fullThreadLoadAllPosts(ft *FullThread, rows []queries.ListThreadPostsRow) {
	ft.AllPosts = make([]Post, len(rows))
	for i, row := range rows {
		ft.AllPosts[i].ID = uint64(*row.PostID)
		ft.AllPosts[i].CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
		ft.AllPosts[i].Author = row.Author.String
		ft.AllPosts[i].Body = *row.Body
	}
}

type DBRepo struct {
	Querier queries.Querier
}

func (r *DBRepo) ListBoards(ctx context.Context) ([]Board, error) {
	dbBoards, err := r.Querier.ListBoards(ctx)
	if err != nil {
		return nil, err
	}

	boards := make([]Board, len(dbBoards))
	for i := range boards {
		boardFromList(&boards[i], dbBoards[i])
	}

	return boards, nil
}

var ErrPageInvalid = errors.New("page invalid, must be 1 or greater")

func (r *DBRepo) ShowBoardListRecenthreads(ctx context.Context, boardSlug string, page int) (*BoardRecentThreads, error) {
	if page < 0 {
		return nil, ErrPageInvalid
	}

	dbBoard, err := r.Querier.ShowBoardFromSlug(ctx, boardSlug)
	if err != nil {
		return nil, err
	}

	dbThreads, err := r.Querier.ListActiveBoardThreads(ctx, queries.ListActiveBoardThreadsParams{
		BoardID: dbBoard.BoardID,
		Limit:   15,
		Offset:  (page - 1) * 15,
	})
	if err != nil {
		return nil, err
	}

	threadIDs := make([]int, len(dbThreads))
	for i := range dbThreads {
		threadIDs[i] = *dbThreads[i].ThreadID
	}

	dbRecentPosts, err := r.Querier.ListThreadRecentPosts(ctx, dbBoard.BoardID, threadIDs)
	if err != nil {
		return nil, err
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

	board := Board{}
	boardFromShow(&board, dbBoard)

	recentThreads := make([]RecentThread, len(dbThreads))
	for i := range dbThreads {
		recentThreadFromActive(&recentThreads[i], dbThreads[i])
		dbRecentPosts, ok := dbRecentPostsByThreads[*dbThreads[i].ThreadID]
		if ok {
			recentThreadLoadRecentPosts(&recentThreads[i], dbRecentPosts)
		} else {
			recentThreadLoadRecentPosts(&recentThreads[i], []queries.ListThreadRecentPostsRow{})
		}
	}

	return &BoardRecentThreads{
		Board:         board,
		RecentThreads: recentThreads,
	}, nil
}

func (r *DBRepo) ShowThread(ctx context.Context, boardSlug string, threadID int) (*BoardFullThread, error) {
	dbBoard, err := r.Querier.ShowBoardFromSlug(ctx, boardSlug)
	if err != nil {
		return nil, err
	}

	dbThread, err := r.Querier.ShowThread(ctx, dbBoard.BoardID, threadID)
	if err != nil {
		return nil, err
	}

	dbPosts, err := r.Querier.ListThreadPosts(ctx, dbBoard.BoardID, threadID)
	if err != nil {
		return nil, err
	}

	board := Board{}
	boardFromShow(&board, dbBoard)

	fullThread := FullThread{}
	fullThreadFromShow(&fullThread, dbThread)
	fullThreadLoadAllPosts(&fullThread, dbPosts)

	return &BoardFullThread{
		Board:      board,
		FullThread: fullThread,
	}, nil
}

var _ CreamyBoard = &DBRepo{}
