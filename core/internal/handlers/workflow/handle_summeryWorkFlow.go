package workflow

//import (
//	"context"
//	"fmt"
//	"net/http"
//	"time"
//
//	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
//	"github.com/jackc/pgx/v5/pgtype"
//	"github.com/rs/zerolog/log"
//	"go.temporal.io/sdk/client"
//	db "stock-agent.io/db/sqlc"
//	"stock-agent.io/internal/types"
//	"stock-agent.io/utils"
//)
//
//func (w *Handler) SummeryWorkflowSchedule(c *gin.Context) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	userID := c.GetString(types.UserIDContextKey)
//
//	scheduleID := fmt.Sprintf("schedule-%s-%s", userID, time.Now().Format("20060102150405"))
//	var req *types.CreateUserScheduleRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	workflowInDb, err := w.store.CreateWorkflowSchedule(ctx, db.CreateWorkflowScheduleParams{
//		ID:           uuid.New(),
//		RelatedTo:    &userID,
//		WorkflowType: req.WorkflowType,
//		ScheduleID:   scheduleID,
//		CreatedAt:    pgtype.Timestamp{},
//	})
//	if err != nil {
//		log.Debug().Err(err).Msgf("Error creating workflows schedule in db")
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	scheduleSpec, err := utils.ParseHumanFriendlySchedule(req.Schedule)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid schedule: %s", err.Error())})
//		return
//	}
//
//	scheduleOptions := client.ScheduleOptions{
//		ID: scheduleID,
//		Spec: client.ScheduleSpec{
//			Calendars: []client.ScheduleCalendarSpec{{
//				Second:     scheduleSpec.Second,
//				Minute:     scheduleSpec.Minute,
//				Hour:       scheduleSpec.Hour,
//				DayOfMonth: scheduleSpec.DayOfMonth,
//				Month:      scheduleSpec.Month,
//				Year:       scheduleSpec.Year,
//				DayOfWeek:  scheduleSpec.DayOfWeek,
//				Comment:    req.Description,
//			}},
//			Intervals:       nil,
//			CronExpressions: nil, // Remove cron expressions
//			Skip:            nil,
//			StartAt:         time.Time{},
//			EndAt:           time.Time{},
//			Jitter:          0,
//			TimeZoneName:    req.TimeZone, // Add timezone support
//		},
//		Memo: map[string]interface{}{
//			"user_id":       userID,
//			"workflow_type": req.WorkflowType,
//			"description":   req.Description,
//			"schedule":      req.Schedule,
//		},
//	}
//
//	switch req.WorkflowType {
//	case "summery-workflows":
//		scheduleOptions.Action = &client.ScheduleWorkflowAction{
//			ID:        workflowInDb.ID.String(),
//			Workflow:  w.workflowManger.StockSummeryWorkflow,
//			TaskQueue: "default",
//		}
//	default:
//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workflows type"})
//		return
//	}
//
//	createdSchedule, err := w.schedulerClient.Create(ctx, scheduleOptions)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusCreated, createdSchedule)
//}
