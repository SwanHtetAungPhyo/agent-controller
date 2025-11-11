-- name: CreateUser :one
INSERT INTO kainos_user (id, clerk_id, first_name, email) VALUES (@id, @clerk_id, @first_name, @email) returning *;

-- name: GetUserByClerkID :one
SELECT * FROM kainos_user
WHERE clerk_id = @clerk_id AND deleted_at IS NULL;

-- name: UpdateUserByClerkID :one
UPDATE kainos_user
SET first_name = @first_name, email = @email, updated_at = NOW()
WHERE clerk_id = @clerk_id AND deleted_at is NULL
returning *;

-- name: SoftDeleteUserByClerkID :one
UPDATE kainos_user
SET deleted_at = NOW()
WHERE clerk_id = @clerk_id AND deleted_at is NULL
returning *;

-- name: GetUserByID :one
SELECT * FROM kainos_user
WHERE id = @id AND deleted_at is NULL;
