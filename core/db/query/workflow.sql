-- name: CreateWorkflowSchedule :one
INSERT INTO workflows_schedules (
    id,
    related_to,
    workflow_type,
    schedule_id,
    created_at
) VALUES (
             $1, $2, $3, $4, $5
         ) RETURNING *;

-- name: GetWorkflowSchedule :one
SELECT * FROM workflows_schedules
WHERE id = $1 LIMIT 1;

-- name: GetWorkflowScheduleByScheduleID :one
SELECT * FROM workflows_schedules
WHERE schedule_id = $1 LIMIT 1;

-- name: ListWorkflowSchedulesByUser :many
SELECT * FROM workflows_schedules
WHERE related_to = $1
ORDER BY created_at DESC;

-- name: ListWorkflowSchedulesByUserAndType :many
SELECT * FROM workflows_schedules
WHERE related_to = $1 AND workflow_type = $2
ORDER BY created_at DESC;

-- name: GetUserWorkflowSchedule :one
SELECT * FROM workflows_schedules
WHERE related_to = $1 AND workflow_type = $2
LIMIT 1;

-- name: UpdateWorkflowSchedule :one
UPDATE workflows_schedules
SET
    workflow_type = $2,
    schedule_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteWorkflowSchedule :exec
DELETE FROM workflows_schedules
WHERE id = $1;

-- name: DeleteWorkflowScheduleByScheduleID :exec
DELETE FROM workflows_schedules
WHERE schedule_id = $1;

-- name: DeleteUserWorkflowSchedules :exec
DELETE FROM workflows_schedules
WHERE related_to = $1;

-- name: DeleteUserWorkflowScheduleByType :exec
DELETE FROM workflows_schedules
WHERE related_to = $1 AND workflow_type = $2;

-- name: ListAllWorkflowSchedules :many
SELECT * FROM workflows_schedules
ORDER BY created_at DESC;

-- name: CountUserWorkflowSchedules :one
SELECT COUNT(*) FROM workflows_schedules
WHERE related_to = $1;

-- name: WorkflowScheduleExists :one
SELECT EXISTS(
    SELECT 1 FROM workflows_schedules
    WHERE related_to = $1 AND workflow_type = $2
);
