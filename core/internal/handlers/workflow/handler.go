package workflow

import (
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

func (w *Handler) RegisterRoutes() {
	// Routes will be registered by the server module
}
