-- name: ListBoards :many
SELECT board_id, slug, title, tagline
FROM boards
ORDER BY slug, board_id
;

-- name: ShowBoardFromID :one
SELECT board_id, slug, title, tagline, (SELECT COUNT(*) FROM threads WHERE threads.board_id = boards.board_id) AS threads
FROM boards
WHERE board_id = pggen.arg('id')
;

-- name: ShowBoardFromSlug :one
SELECT board_id, slug, title, tagline, (SELECT COUNT(*) FROM threads WHERE threads.board_id = boards.board_id) AS threads
FROM boards
WHERE slug = pggen.arg('slug')
;

-- name: ListActiveBoardThreads :many
SELECT threads.thread_id, threads.created_at, threads.bumped_at, posts.subject, posts.author, posts.body
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
SELECT threads.thread_id, recent_posts.post_id, recent_posts.created_at, recent_posts.subject, recent_posts.author, recent_posts.body
FROM threads
JOIN LATERAL (
  SELECT post_id, created_at, subject, author, body
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
SELECT threads.thread_id, threads.created_at, threads.bumped_at, posts.subject, posts.author, posts.body
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
SELECT post_id, created_at, subject, author, body
FROM posts
WHERE posts.board_id = pggen.arg('board_id')
AND posts.thread_id = pggen.arg('thread_id')
AND posts.post_id != pggen.arg('thread_id') -- ignore the thread post
ORDER BY posts.post_id
;

-- name: ListPostFiles :many
SELECT 
post_id, idx
, path, extension, mimetype, bytes, original_name
, thumb_path, thumb_extension, thumb_mimetype, thumb_bytes
FROM files
WHERE board_id = pggen.arg('board_id')
AND post_id = ANY (pggen.arg('post_ids')::BIGINT[])
ORDER BY post_id, idx
;

-- name: ListThreadFiles :many
SELECT
post_id, idx
, path, extension, mimetype, bytes, original_name
, thumb_path, thumb_extension, thumb_mimetype, thumb_bytes
FROM files
WHERE board_id = pggen.arg('board_id')
AND thread_id = pggen.arg('thread_id')
ORDER BY post_id, idx
;

-- name: ShowFile :one
SELECT extension, path, mimetype, bytes
FROM files
WHERE board_id = pggen.arg('board_id')
AND thread_id = pggen.arg('thread_id')
AND post_id = pggen.arg('post_id')
AND idx = pggen.arg('idx')
;

-- name: ShowFileThumb :one
SELECT thumb_extension, thumb_path, thumb_mimetype, thumb_bytes
FROM files
WHERE board_id = pggen.arg('board_id')
AND thread_id = pggen.arg('thread_id')
AND post_id = pggen.arg('post_id')
AND idx = pggen.arg('idx')
AND thumb_path IS NOT NULL
;

-- name: SubmitThread :one
WITH
thread AS (
  INSERT INTO threads (board_id, thread_id) VALUES
    (pggen.arg('board_id'), board_post_seq_nextval(pggen.arg('board_id')))
  RETURNING thread_id
),
post AS (
  INSERT INTO posts (board_id, thread_id, post_id, subject, author, body) VALUES
    (pggen.arg('board_id'), (SELECT thread_id FROM thread), (SELECT thread_id FROM thread), pggen.arg('subject'), pggen.arg('author'), pggen.arg('body'))
  RETURNING post_id
),
files_input AS (
  SELECT pggen.arg('board_id') AS board_id, (SELECT thread_id FROM thread) AS thread_id, (SELECT thread_id FROM thread) AS post_id, idx, path, extension, mimetype, bytes, original_name, thumb_path, thumb_extension, thumb_mimetype, thumb_bytes
  FROM unnest(pggen.arg('partial_files')::partial_file[])
),
files AS (
  INSERT INTO files (board_id, thread_id, post_id, idx, path, extension, mimetype, bytes, original_name, thumb_path, thumb_extension, thumb_mimetype, thumb_bytes)
  SELECT *
  FROM files_input
)
SELECT thread_id FROM thread
;

-- name: SubmitPost :one
WITH
post AS (
  INSERT INTO posts (board_id, thread_id, post_id, subject, author, body) VALUES
    (pggen.arg('board_id'), pggen.arg('thread_id'), board_post_seq_nextval(pggen.arg('board_id')), pggen.arg('subject'), pggen.arg('author'), pggen.arg('body'))
  RETURNING post_id
),
files_input AS (
  SELECT pggen.arg('board_id') AS board_id, pggen.arg('thread_id'), (SELECT post_id FROM post) AS post_id, idx, path, extension, mimetype, bytes, original_name, thumb_path, thumb_extension, thumb_mimetype, thumb_bytes
  FROM unnest(pggen.arg('partial_files')::partial_file[])
),
files AS (
  INSERT INTO files (board_id, thread_id, post_id, idx, path, extension, mimetype, bytes, original_name, thumb_path, thumb_extension, thumb_mimetype, thumb_bytes)
  SELECT *
  FROM files_input
)
SELECT post_id FROM post
;

-- the above SubmitPost query fails when no files are submitted, not sure how to hack around it, so just using diff query instead
-- name: SubmitPostNoFiles :one
INSERT INTO posts (board_id, thread_id, post_id, subject, author, body) VALUES
    (pggen.arg('board_id'), pggen.arg('thread_id'), board_post_seq_nextval(pggen.arg('board_id')), pggen.arg('subject'), pggen.arg('author'), pggen.arg('body'))
  RETURNING post_id
;

-- name: BumpThread :exec
UPDATE threads
SET bumped_at = NOW()
WHERE board_id = pggen.arg('board_id')
AND thread_id = pggen.arg('thread_id')
;
