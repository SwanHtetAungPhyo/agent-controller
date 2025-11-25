package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// ExecuteMastraWorkflow - Temporal workflow that calls Mastra AI
func (m *Manager) ExecuteMastraWorkflow(ctx workflow.Context, userWorkflowID, workflowID string) error {
	// Set activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Call Mastra API activity
	var result string
	err := workflow.ExecuteActivity(ctx, m.CallMastraAPI, userWorkflowID, workflowID).Get(ctx, &result)
	if err != nil {
		return fmt.Errorf("failed to call Mastra API: %w", err)
	}

	// Store result activity
	err = workflow.ExecuteActivity(ctx, m.StoreWorkflowResult, userWorkflowID, result).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to store result: %w", err)
	}

	return nil
}

// CallMastraAPI - Activity that makes HTTP call to Mastra (Mock for now)
func (m *Manager) CallMastraAPI(ctx context.Context, userWorkflowID, workflowID string) (string, error) {
	// Mock Mastra API call
	mastraURL := "https://api.mastra.ai/workflows/execute" // Example URL

	log.Info().
		Str("user_workflow_id", userWorkflowID).
		Str("workflow_id", workflowID).
		Str("mastra_url", mastraURL).
		Msg("Calling Mastra AI API (MOCK)")

	// TODO: Replace with actual HTTP call
	// client := &http.Client{Timeout: 30 * time.Second}
	// payload := map[string]interface{}{
	//     "workflow_id": workflowID,
	//     "user_workflow_id": userWorkflowID,
	//     "timestamp": time.Now(),
	// }
	// resp, err := client.Post(mastraURL, "application/json", bytes.NewBuffer(jsonPayload))

	// For now, return mock response
	mockResult := fmt.Sprintf("Mock result from Mastra AI for workflow %s at %s", workflowID, time.Now().Format(time.RFC3339))

	return mockResult, nil
}

// StoreWorkflowResult - Activity to store workflow execution result
func (m *Manager) StoreWorkflowResult(ctx context.Context, userWorkflowID, result string) error {
	log.Info().
		Str("user_workflow_id", userWorkflowID).
		Str("result", result).
		Msg("Storing workflow result (MOCK)")

	// TODO: Store in database
	// m.store.CreateWorkflowExecution(ctx, db.CreateWorkflowExecutionParams{
	//     UserWorkflowID: userWorkflowID,
	//     Result: result,
	//     ExecutedAt: time.Now(),
	// })

	return nil
}
