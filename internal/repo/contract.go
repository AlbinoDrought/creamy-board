package repo

import "context"

type Board struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Tagline string `json:"tagline"`
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

type FullThread struct {
	Thread   Thread `json:"thread"`
	MainPost Post   `json:"main_post"`
	AllPosts []Post `json:"all_posts"`
}

type CreamyBoard interface {
	ListBoards(ctx context.Context) ([]Board, error)
	ShowBoard(ctx context.Context, boardSlug string) (*Board, error)
	ListBoardThreads(ctx context.Context, boardSlug string, page int) ([]RecentThread, error)
	ShowThread(ctx context.Context, boardSlug string, threadID int) (*FullThread, error)
}
