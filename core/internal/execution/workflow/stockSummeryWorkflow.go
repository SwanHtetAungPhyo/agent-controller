package workflow

import (
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/workflow"
)

func (w *Manager) StockSummeryWorkflow(ctx workflow.Context) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Stock Summery workflow started")
	logger.Info("Starting workflow")
	log.Debug().Msg("Starting workflow")

}
