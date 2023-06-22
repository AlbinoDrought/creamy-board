-- name: ListBoards :many
SELECT board_id, slug, title, tagline
FROM boards
ORDER BY slug, board_id
;

-- name: ShowBoardFromID :one
SELECT board_id, slug, title, tagline
FROM boards
WHERE board_id = pggen.arg('id')
;

-- name: ShowBoardFromSlug :one
SELECT board_id, slug, title, tagline
FROM boards
WHERE slug = pggen.arg('slug')
;

-- name: ListActiveBoardThreads :many
SELECT threads.thread_id, threads.created_at, threads.bumped_at, posts.author, posts.body
FROM threads
-- join the thread post:
INNER JOIN posts
ON posts.board_id = threads.board_id
AND posts.thread_id = threads.thread_id
AND posts.post_id = threads.thread_id
WHERE threads.board_id = pggen.arg('board_id')
ORDER BY threads.bumped_at DESC
LIMIT pggen.arg('limit')
OFFSET pggen.arg('offset') 
;

-- name: ListThreadRecentPosts :many
SELECT threads.thread_id, recent_posts.post_id, recent_posts.created_at, recent_posts.author, recent_posts.body
FROM threads
JOIN LATERAL (
  SELECT post_id, created_at, author, body
  FROM posts
  WHERE posts.board_id = threads.board_id
  AND posts.thread_id = threads.thread_id
  AND posts.post_id != threads.thread_id -- ignore the thread post
  ORDER BY posts.post_id
  LIMIT 5
) recent_posts ON TRUE
WHERE threads.board_id = pggen.arg('board_id')
AND threads.thread_id = ANY (pggen.arg('thread_ids')::BIGINT[])
;

-- name: ShowThread :one
SELECT threads.thread_id, threads.created_at, threads.bumped_at, posts.author, posts.body
FROM threads
-- join the thread post:
INNER JOIN posts
ON posts.board_id = threads.board_id
AND posts.thread_id = threads.thread_id
AND posts.post_id = threads.thread_id
WHERE threads.board_id = pggen.arg('board_id')
AND threads.thread_id = pggen.arg('thread_id')
;

-- name: ListThreadPosts :many
SELECT post_id, created_at, author, body
FROM posts
WHERE posts.board_id = pggen.arg('board_id')
AND posts.thread_id = pggen.arg('thread_id')
AND posts.post_id != pggen.arg('thread_id') -- ignore the thread post
ORDER BY posts.post_id
;
