-- -- name: CreateWorkflowSchedule :one
-- INSERT INTO workflows_schedules (
--     id,
--     related_to,
--     workflow_type,
--     schedule_id,
--     created_at
-- ) VALUES (
--              $1, $2, $3, $4, $5
--          ) RETURNING *;
--
-- -- name: GetWorkflowSchedule :one
-- SELECT * FROM workflows_schedules
-- WHERE id = $1 LIMIT 1;
--
-- -- name: GetWorkflowScheduleByScheduleID :one
-- SELECT * FROM workflows_schedules
-- WHERE schedule_id = $1 LIMIT 1;
--
-- -- name: ListWorkflowSchedulesByUser :many
-- SELECT * FROM workflows_schedules
-- WHERE related_to = $1
-- ORDER BY created_at DESC;
--
-- -- name: ListWorkflowSchedulesByUserAndType :many
-- SELECT * FROM workflows_schedules
-- WHERE related_to = $1 AND workflow_type = $2
-- ORDER BY created_at DESC;
--
-- -- name: GetUserWorkflowSchedule :one
-- SELECT * FROM workflows_schedules
-- WHERE related_to = $1 AND workflow_type = $2
-- LIMIT 1;
--
-- -- name: UpdateWorkflowSchedule :one
-- UPDATE workflows_schedules
-- SET
--     workflow_type = $2,
--     schedule_id = $3
-- WHERE id = $1
-- RETURNING *;
--
-- -- name: DeleteWorkflowSchedule :exec
-- DELETE FROM workflows_schedules
-- WHERE id = $1;
--
-- -- name: DeleteWorkflowScheduleByScheduleID :exec
-- DELETE FROM workflows_schedules
-- WHERE schedule_id = $1;
--
-- -- name: DeleteUserWorkflowSchedules :exec
-- DELETE FROM workflows_schedules
-- WHERE related_to = $1;
--
-- -- name: DeleteUserWorkflowScheduleByType :exec
-- DELETE FROM workflows_schedules
-- WHERE related_to = $1 AND workflow_type = $2;
--
-- -- name: ListAllWorkflowSchedules :many
-- SELECT * FROM workflows_schedules
-- ORDER BY created_at DESC;
--
-- -- name: CountUserWorkflowSchedules :one
-- SELECT COUNT(*) FROM workflows_schedules
-- WHERE related_to = $1;
--
-- -- name: WorkflowScheduleExists :one
-- SELECT EXISTS(
--     SELECT 1 FROM workflows_schedules
--     WHERE related_to = $1 AND workflow_type = $2
-- );

-- name: CreateWorkflow :one
INSERT INTO kainos_workflow (id, workflow_name, workflow_description, price) VALUES (@id,@workflow_name, @workflow_description, @price) returning *;

-- name: CreateUserWorkflow :one
INSERT INTO kainos_user_workflow (id, workflow_id, customer_id, meta_data, status) VALUES (@id,@workflow_id, @customer_id, @meta_data, @status) returning *;

-- name: GetWorkflow :many
SELECT * from kainos_workflow;

-- name: GetUserWorkflow :many
SELECT workflow_id, workflow_name, meta_data, cron_time, status, kainos_user_workflow.created_at, kainos_user_workflow.updated_at
from kainos_user_workflow
         join kainos_workflow
              on kainos_user_workflow.workflow_id = kainos_workflow.id
         join kainos_user
              on kainos_user_workflow.customer_id = kainos_user.id;