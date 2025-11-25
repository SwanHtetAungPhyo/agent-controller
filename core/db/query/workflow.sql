-- name: CreateWorkflow :one
INSERT INTO kainos_workflow (id, workflow_name, workflow_description, price) VALUES (@id, @workflow_name, @workflow_description, @price) returning *;

-- name: CreateUserWorkflow :one
INSERT INTO kainos_user_workflow (id, workflow_id, customer_id, meta_data, status) VALUES (@id, @workflow_id, @customer_id, @meta_data, @status) returning *;

-- name: GetWorkflow :many
SELECT * from kainos_workflow
WHERE deleted_at is NULL;

-- -- name: GetUserWorkflow :many
-- SELECT workflow_id, workflow_name, meta_data, cron_time, status, kainos_user_workflow.created_at, kainos_user_workflow.updated_at
-- from kainos_user_workflow
-- join kainos_workflow on kainos_user_workflow.workflow_id = kainos_workflow.id
-- join kainos_user on kainos_user_workflow.customer_id = kainos_user.id;

-- name: GetUserWorkflowsByClerkID :many
SELECT uw.id, uw.workflow_id, uw.customer_id, uw.meta_data, uw.cron_time, uw.status, uw.created_at,
       uw.updated_at, w.workflow_name, w.workflow_description, w.price
FROM kainos_user_workflow uw
JOIN kainos_workflow w ON uw.workflow_id = w.id
JOIN kainos_user u ON uw.customer_id = u.id
WHERE u.clerk_id = @clerk_id AND u.deleted_at is NULL;

-- name: UpdateUserWorkflowStatus :one
UPDATE kainos_user_workflow
SET status = @status, updated_at = NOW()
WHERE id = @id
returning *;

-- name: UpdateUserWorkflowSchedule :one
UPDATE kainos_user_workflow
SET
    cron_time = @cron_time,
    status = @status,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: GetUserWorkflowByID :one
SELECT
    uw.id,
    uw.workflow_id,
    uw.customer_id,
    uw.meta_data,
    uw.cron_time,
    uw.status,
    uw.created_at,
    uw.updated_at,
    w.workflow_name,
    w.workflow_description,
    w.price
FROM kainos_user_workflow uw
         JOIN kainos_workflow w ON uw.workflow_id = w.id
WHERE uw.id = @id;
