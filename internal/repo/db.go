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

func fileFromListPost(f *File, row queries.ListPostFilesRow) {
	f.Index = int(*row.Idx)
	f.InternalPath = row.Path.String
	f.Extension = row.Extension.String
	f.MimeType = row.Mimetype.String
	f.Bytes = int(*row.Bytes)
	f.OriginalName = row.OriginalName.String
}

func fileFromListThread(f *File, row queries.ListThreadFilesRow) {
	f.Index = int(*row.Idx)
	f.InternalPath = row.Path.String
	f.Extension = row.Extension.String
	f.MimeType = row.Mimetype.String
	f.Bytes = int(*row.Bytes)
	f.OriginalName = row.OriginalName.String
}

func postFromRecentPost(p *Post, row queries.ListThreadRecentPostsRow) {
	p.ID = uint64(*row.PostID)
	p.CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
	p.Author = row.Author.String
	p.Body = *row.Body
}

func postFromThreadPost(p *Post, row queries.ListThreadPostsRow) {
	p.ID = uint64(*row.PostID)
	p.CreatedAt = row.CreatedAt.Time.Format(time.RFC3339)
	p.Author = row.Author.String
	p.Body = *row.Body
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

func (r *DBRepo) ShowBoardListRecentThreads(ctx context.Context, boardSlug string, page int) (*BoardRecentThreads, error) {
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

	// query recent posts for all threads
	threadIDs := make([]int, len(dbThreads))
	for i := range dbThreads {
		threadIDs[i] = *dbThreads[i].ThreadID
	}

	dbRecentPosts, err := r.Querier.ListThreadRecentPosts(ctx, dbBoard.BoardID, threadIDs)
	if err != nil {
		return nil, err
	}

	// query files for all posts
	postIDs := make([]int, len(dbRecentPosts)+1)
	for i, dbRecentPost := range dbRecentPosts {
		postIDs[i] = *dbRecentPost.PostID
	}

	threadsAndPosts := len(dbThreads) + len(dbRecentPosts)

	dbFiles, err := r.Querier.ListPostFiles(ctx, dbBoard.BoardID, append(threadIDs, postIDs...))
	if err != nil {
		return nil, err
	}

	// map of post ID -> post files
	filesByPosts := make(map[int][]File, threadsAndPosts)
	for _, dbFile := range dbFiles {
		file := File{}
		fileFromListPost(&file, dbFile)

		filesByPost, ok := filesByPosts[*dbFile.PostID]
		if !ok {
			filesByPost = []File{}
		}
		filesByPost = append(filesByPost, file)
		filesByPosts[*dbFile.PostID] = filesByPost
	}

	// map of thread ID -> recent posts
	recentPostsByThreads := make(map[int][]Post, len(dbThreads))
	for _, dbRecentPost := range dbRecentPosts {
		post := Post{}
		postFromRecentPost(&post, dbRecentPost)

		files, ok := filesByPosts[int(post.ID)]
		if ok {
			post.Files = files
		} else {
			post.Files = []File{}
		}

		recentPostsByThread, ok := recentPostsByThreads[*dbRecentPost.ThreadID]
		if !ok {
			recentPostsByThread = []Post{}
		}
		recentPostsByThread = append(recentPostsByThread, post)
		recentPostsByThreads[*dbRecentPost.ThreadID] = recentPostsByThread
	}

	board := Board{}
	boardFromShow(&board, dbBoard)

	recentThreads := make([]RecentThread, len(dbThreads))
	for i := range dbThreads {
		recentThreadFromActive(&recentThreads[i], dbThreads[i])

		files, ok := filesByPosts[int(*dbThreads[i].ThreadID)]
		if ok {
			recentThreads[i].MainPost.Files = files
		} else {
			recentThreads[i].MainPost.Files = []File{}
		}

		recentPosts, ok := recentPostsByThreads[*dbThreads[i].ThreadID]
		if ok {
			recentThreads[i].RecentPosts = recentPosts
		} else {
			recentThreads[i].RecentPosts = []Post{}
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

	threadsAndPosts := 1 + len(dbPosts)

	dbFiles, err := r.Querier.ListThreadFiles(ctx, dbBoard.BoardID, threadID)
	if err != nil {
		return nil, err
	}

	// map of post ID -> post files
	filesByPosts := make(map[int][]File, threadsAndPosts)
	for _, dbFile := range dbFiles {
		file := File{}
		fileFromListThread(&file, dbFile)

		filesByPost, ok := filesByPosts[*dbFile.PostID]
		if !ok {
			filesByPost = []File{}
		}
		filesByPost = append(filesByPost, file)
		filesByPosts[*dbFile.PostID] = filesByPost
	}

	posts := make([]Post, len(dbPosts))
	for i := range dbPosts {
		postFromThreadPost(&posts[i], dbPosts[i])

		files, ok := filesByPosts[int(posts[i].ID)]
		if ok {
			posts[i].Files = files
		} else {
			posts[i].Files = []File{}
		}
	}

	board := Board{}
	boardFromShow(&board, dbBoard)

	fullThread := FullThread{}
	fullThreadFromShow(&fullThread, dbThread)
	fullThread.AllPosts = posts

	files, ok := filesByPosts[int(fullThread.MainPost.ID)]
	if ok {
		fullThread.MainPost.Files = files
	} else {
		fullThread.MainPost.Files = []File{}
	}

	return &BoardFullThread{
		Board:      board,
		FullThread: fullThread,
	}, nil
}

var _ CreamyBoard = &DBRepo{}
