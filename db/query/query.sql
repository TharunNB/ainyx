-- name: CreateUser :one
-- Creates a new user record, returns created user
INSERT INTO users (name, dob)
VALUES ($1, $2)
RETURNING id, name, dob;

-- name: GetUserByID :one
-- Fetches a single user by key

SELECT id, name, dob
FROM users
WHERE id = $1;

-- name: UpdateUser :one
-- Updates a user's name and dob, updated row.

UPDATE users
SET name = $1,
    dob = $2
WHERE id = $3
RETURNING id, name, dob;

-- name: DeleteUser :exec
-- Deletes a user by the key

DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
-- Returns list of users ordered by ID

SELECT id, name, dob
FROM users 
ORDER BY id ASC 
LIMIT $1
OFFSET $2;

-- name: CountUsers :one
-- Returns the total number of users 

SELECT COUNT(*) FROM users;