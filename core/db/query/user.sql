-- name: CreateUser :one
INSERT INTO kainos_user (clerk_id, first_name, email) VALUES (@clerk_id, @first_name, @email) returning *;

-- name: CreateWorkflow :one
INSERT INTO kainos_workflow (workflow_name, workflow_description, price) VALUES (@workflow_name, @workflow_description, @price) returning *;

-- name: CreateUserWorkflow :one
INSERT INTO kainos_user_workflow (workflow_id, customer_id, meta_data, status) VALUES (@workflow_id, @customer_id, @meta_data, @status) returning *;

-- name: CreateSystemAnalysis :one
INSERT INTO system_defined_analysis (analysis_type, description) VALUES (@analysis_type, @description) returning *;

-- name: CreateUserAnalysis :one
INSERT INTO kainos_user_analysis (description, s3_url, customer_id) VALUES (@description, @s3_url, @customer_id) returning *;

-- name: GetWorkflow :many
SELECT * from kainos_workflow;

-- name: GetUserWorkflow :many
SELECT workflow_id, workflow_name, meta_data, cron_time, status, kainos_user_workflow.created_at, kainos_user_workflow.updated_at
        from kainos_user_workflow
        join kainos_workflow
            on kainos_user_workflow.workflow_id = kainos_workflow.id
        join kainos_user
        on kainos_user_workflow.customer_id = kainos_user.id;
