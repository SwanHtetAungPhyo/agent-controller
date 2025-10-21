-- name: CreateUser :one
INSERT INTO kainos_user (id, clerk_id, first_name, email) VALUES (@id,@clerk_id, @first_name, @email) returning *;

