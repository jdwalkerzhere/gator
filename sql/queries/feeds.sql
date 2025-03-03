-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT
    f.name AS feed_name,
    f.url,
    u.name AS user_name
FROM feeds f
INNER JOIN users u ON f.user_id = u.id;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET
    (updated_at, last_fetched_at) = (NOW(), NOW())
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at NULLS FIRST, updated_at
LIMIT 1;
