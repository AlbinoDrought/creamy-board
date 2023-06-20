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
