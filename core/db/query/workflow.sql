-- name: CreateWorkflow :one
INSERT INTO kainos_workflow (id, workflow_name, workflow_description, price) VALUES (@id, @workflow_name, @workflow_description, @price) returning *;

-- name: CreateUserWorkflow :one
INSERT INTO kainos_user_workflow (id, workflow_id, customer_id, meta_data, status) VALUES (@id, @workflow_id, @customer_id, @meta_data, @status) returning *;

-- name: GetWorkflow :many
SELECT * from kainos_workflow
WHERE deleted_at is NULL;

-- name: GetUserWorkflow :many
SELECT workflow_id, workflow_name, meta_data, cron_time, status, kainos_user_workflow.created_at, kainos_user_workflow.updated_at
from kainos_user_workflow
join kainos_workflow on kainos_user_workflow.workflow_id = kainos_workflow.id
join kainos_user on kainos_user_workflow.customer_id = kainos_user.id;

-- name: UpdateUserWorkflowStatus :one
UPDATE kainos_user_workflow
SET status = @status, updated_at = NOW()
WHERE id = @id
returning *;
