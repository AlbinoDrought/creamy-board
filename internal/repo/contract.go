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

type File struct {
	Index        int    `json:"index"`
	Extension    string `json:"extension"`
	MimeType     string `json:"mime_type"`
	Bytes        int    `json:"bytes"`
	OriginalName string `json:"original_name"`

	InternalPath string `json:"-"`
}

type Post struct {
	ID        uint64 `json:"id"`
	CreatedAt string `json:"created_at"`
	Subject   string `json:"subject"`
	Author    string `json:"author"`
	Body      string `json:"body"`
	Files     []File `json:"files"`
}

type RecentThread struct {
	Thread      Thread `json:"thread"`
	MainPost    Post   `json:"main_post"`
	RecentPosts []Post `json:"recent_posts"`
}

type BoardRecentThreads struct {
	Board         Board          `json:"board"`
	RecentThreads []RecentThread `json:"recent_threads"`
}

type FullThread struct {
	Thread   Thread `json:"thread"`
	MainPost Post   `json:"main_post"`
	AllPosts []Post `json:"all_posts"`
}

type BoardFullThread struct {
	Board      Board      `json:"board"`
	FullThread FullThread `json:"full_thread"`
}

type CreamyBoard interface {
	ListBoards(ctx context.Context) ([]Board, error)
	ShowBoardListRecentThreads(ctx context.Context, boardSlug string, page int) (*BoardRecentThreads, error)
	ShowThread(ctx context.Context, boardSlug string, threadID int) (*BoardFullThread, error)
}
