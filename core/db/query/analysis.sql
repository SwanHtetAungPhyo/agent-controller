-- name: CreateSystemAnalysis :one
INSERT INTO system_defined_analysis (id, analysis_type, description) VALUES (@id,@analysis_type, @description) returning *;

-- name: CreateUserAnalysis :one
INSERT INTO kainos_user_analysis (id, description, s3_url, customer_id) VALUES (@id,@description, @s3_url, @customer_id) returning *;

-- name: GetSystemAnalysis :many
SELECT * from system_defined_analysis;

-- name: GetUserAnalysis :many
SELECT kainos_user_analysis.id, description, s3_url,kainos_user_analysis.created_at
    from kainos_user_analysis
    join kainos_user
    on kainos_user_analysis.customer_id = kainos_user.id;