package types

type CreateUserScheduleRequest struct {
	WorkflowType string `json:"workflow_type" binding:"required"`
	Schedule     string `json:"schedule" binding:"required"`
	Description  string `json:"description"`
	TimeZone     string `json:"time_zone" default:"America/New_York"`
}
