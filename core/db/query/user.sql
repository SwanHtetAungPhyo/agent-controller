-- name: CreateAppUser :one
INSERT  INTO app_users (clerk_id) VALUES  (clerk_id) RETURNING  *;