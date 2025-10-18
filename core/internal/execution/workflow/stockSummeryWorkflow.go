package workflow

import (
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/workflow"
)

func (w *Manager) StockSummeryWorkflow(ctx workflow.Context) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Stock Summery workflows started")
	logger.Info("Starting workflows")
	log.Debug().Msg("Starting workflows")

}
