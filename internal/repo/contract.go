package repo

import (
	"context"
)

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

	ThumbExtension string `json:"thumb_extension"`
	ThumbMimeType  string `json:"thumb_mime_type"`
	ThumbBytes     int    `json:"thumb_bytes"`

	InternalPath      string `json:"-"`
	ThumbInternalPath string `json:"-"`
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
	Pages         int
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

type SubmitPostFile struct {
	Extension    string
	MimeType     string
	Bytes        int
	OriginalName string
	InternalPath string

	ThumbExtension    string
	ThumbMimeType     string
	ThumbBytes        int
	ThumbInternalPath string
}

type SubmitPost struct {
	Subject string
	Author  string
	Body    string
	Files   []SubmitPostFile
}

type CreamyBoard interface {
	ListBoards(ctx context.Context) ([]Board, error)
	ShowBoardListRecentThreads(ctx context.Context, boardSlug string, page int) (*BoardRecentThreads, error)
	ShowThread(ctx context.Context, boardSlug string, threadID int) (*BoardFullThread, error)

	SubmitThread(ctx context.Context, boardSlug string, req SubmitPost) (int, error)
	SubmitThreadPost(ctx context.Context, boardSlug string, threadID int, req SubmitPost) (int, error)
}
