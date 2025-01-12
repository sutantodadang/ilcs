-- name: InsertTodo :one
INSERT INTO todo (id, title, description, due_date) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ListTodo :many
WITH filtered_todo AS (
    SELECT *
    FROM todo
    WHERE 
        (sqlc.arg(status)::text IS NULL OR status = sqlc.arg(status)::todo_status) AND
        (sqlc.arg(search)::text IS NULL OR 
            (title ILIKE '%' || sqlc.arg(search) || '%' OR 
             description ILIKE '%' || sqlc.arg(search) || '%'))
)
SELECT 
    id,
    title,
    description,
    status,
    due_date
FROM filtered_todo
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_val)::integer
OFFSET (sqlc.arg(page)::integer - 1) * sqlc.arg(limit_val)::integer;

-- name: CountTodo :one
SELECT COUNT(*) 
FROM todo
WHERE 
    (sqlc.arg(status)::text IS NULL OR status = sqlc.arg(status)::todo_status) AND
    (sqlc.arg(search)::text IS NULL OR 
        (title ILIKE '%' || sqlc.arg(search) || '%' OR 
         description ILIKE '%' || sqlc.arg(search) || '%'));


-- name: UpdateTodo :one
UPDATE todo 
SET 
    title = $2,
    description = $3,
    status = $4,
    due_date = $5,
    updated_at = NOW()
WHERE id = $1 
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todo WHERE id = $1;

-- name: GetTodoById :one
SELECT 
    id,
    title,
    description,
    status,
    due_date
FROM todo
WHERE id = $1;
