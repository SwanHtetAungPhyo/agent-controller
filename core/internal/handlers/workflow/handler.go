package workflow

import (
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
	db "stock-agent.io/db/sqlc"
	"stock-agent.io/internal/execution/workflow"
	"stock-agent.io/internal/middleware"
)

type Handler struct {
	workflowManger   *workflow.Manager
	schedulerClient  client.ScheduleClient
	app              *gin.Engine
	middleWareManger *middleware.Manager
	store            db.Store
}

func NewHandler(
	app *gin.Engine, // Add this parameter
	workflowManger *workflow.Manager,
	scheduleClient client.ScheduleClient,
	middleWareManager *middleware.Manager,
	store db.Store,
) *Handler {

	return &Handler{
		app:              app, // Initialize the app field
		workflowManger:   workflowManger,
		schedulerClient:  scheduleClient,
		middleWareManger: middleWareManager,
		store:            store,
	}
}

func (w *Handler) RegisterRoutes() {
	workflowApi := w.app.Group("/workflow")
	workflowApi.POST("/", w.middleWareManger.AuthMiddleware(), w.SummeryWorkflowSchedule)
}
