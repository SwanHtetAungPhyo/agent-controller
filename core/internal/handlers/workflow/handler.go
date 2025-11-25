package workflow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	db "stock-agent.io/db/sqlc"
	"stock-agent.io/internal/execution/workflow"
	"stock-agent.io/internal/middleware"
)

type Handler struct {
	workflowManger   *workflow.Manager
	schedulerClient  client.ScheduleClient
	middleWareManger *middleware.Manager
	store            db.Store
}

func NewHandler(
	workflowManger *workflow.Manager,
	scheduleClient client.ScheduleClient,
	middleWareManager *middleware.Manager,
	store db.Store,
) *Handler {
	return &Handler{
		workflowManger:   workflowManger,
		schedulerClient:  scheduleClient,
		middleWareManger: middleWareManager,
		store:            store,
	}
}

func (w *Handler) RegisterRoutes(router *gin.Engine) {
	// Routes will be registered by the server module
	api := router.Group("/api/v1/workflows")
	{
		// User workflow management
		api.GET("/my-workflows", w.GetMyWorkflows)
		api.PATCH("/:id/schedule", w.UpdateWorkflowSchedule)
		api.PATCH("/:id/status", w.UpdateWorkflowStatus)
	}
}

func (w *Handler) GetMyWorkflows(c *gin.Context) {
	clerkID := c.Query("clerk_id")
	if clerkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clerk_id is required"})
		return
	}

	workflows, err := w.store.GetUserWorkflowsByClerkID(c.Request.Context(), clerkID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user workflows")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user workflows"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"workflows": workflows,
		"count":     len(workflows),
	})
}

func (w *Handler) UpdateWorkflowSchedule(c *gin.Context) {
	workflowID := c.Param("id")

	var req struct {
		CronTime string `json:"cron_time" binding:"required"` // e.g., "0 9 * * *"
		Status   string `json:"status" binding:"required"`    // "ON" or "OFF"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	if req.Status != "ON" && req.Status != "OFF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be ON or OFF"})
		return
	}

	// Parse workflow ID
	id, err := uuid.Parse(workflowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// Update in database
	workflow, err := w.store.UpdateUserWorkflowSchedule(c.Request.Context(), db.UpdateUserWorkflowScheduleParams{
		ID:       id,
		CronTime: &req.CronTime,
		Status:   &req.Status,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update workflow")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workflow"})
		return
	}

	// Handle Temporal scheduling
	if req.Status == "ON" {
		if err := w.scheduleWorkflow(c.Request.Context(), workflow); err != nil {
			log.Error().Err(err).Msg("Failed to schedule workflow in Temporal")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule workflow"})
			return
		}
		log.Info().
			Str("workflow_id", workflow.ID.String()).
			Str("cron_time", *workflow.CronTime).
			Msg("Workflow scheduled in Temporal")
	} else {
		if err := w.deleteSchedule(c.Request.Context(), workflow.ID.String()); err != nil {
			log.Error().Err(err).Msg("Failed to delete schedule")
		}
		log.Info().
			Str("workflow_id", workflow.ID.String()).
			Msg("Workflow schedule deleted from Temporal")
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Workflow schedule updated successfully",
		"workflow": workflow,
	})
}

// UpdateWorkflowStatus - Just turn workflow ON/OFF without changing schedule
func (w *Handler) UpdateWorkflowStatus(c *gin.Context) {
	workflowID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"` // "ON" or "OFF"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	if req.Status != "ON" && req.Status != "OFF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be ON or OFF"})
		return
	}

	// Parse workflow ID
	id, err := uuid.Parse(workflowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// Update status in database
	workflow, err := w.store.UpdateUserWorkflowStatus(c.Request.Context(), db.UpdateUserWorkflowStatusParams{
		ID:     id,
		Status: &req.Status,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update workflow status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	// Handle Temporal scheduling
	if req.Status == "ON" && workflow.CronTime != nil {
		if err := w.scheduleWorkflow(c.Request.Context(), workflow); err != nil {
			log.Error().Err(err).Msg("Failed to schedule workflow")
		}
	} else if req.Status == "OFF" {
		if err := w.deleteSchedule(c.Request.Context(), workflow.ID.String()); err != nil {
			log.Error().Err(err).Msg("Failed to delete schedule")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Workflow status updated successfully",
		"workflow": workflow,
	})
}

// scheduleWorkflow - Create Temporal schedule
func (w *Handler) scheduleWorkflow(ctx context.Context, workflow db.KainosUserWorkflow) error {
	if workflow.CronTime == nil {
		return fmt.Errorf("cron_time is required")
	}

	scheduleID := fmt.Sprintf("workflow-%s", workflow.ID.String())

	// Create Temporal schedule
	_, err := w.schedulerClient.Create(ctx, client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			CronExpressions: []string{*workflow.CronTime},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflow.ID.String(),
			Workflow:  w.workflowManger.ExecuteMastraWorkflow, // We'll create this next
			TaskQueue: "default",
			Args: []interface{}{
				workflow.ID.String(),
				workflow.WorkflowID.String(),
			},
		},
	})

	return err
}

// deleteSchedule - Delete Temporal schedule
func (w *Handler) deleteSchedule(ctx context.Context, workflowID string) error {
	scheduleID := fmt.Sprintf("workflow-%s", workflowID)

	handle := w.schedulerClient.GetHandle(ctx, scheduleID)
	return handle.Delete(ctx)
}
