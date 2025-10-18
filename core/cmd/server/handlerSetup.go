package server

import (
	"stock-agent.io/internal/handlers/users"
	"stock-agent.io/internal/handlers/workflow"
)

func (s *HttpServer) handlerSetup() {
	s.handlers = handlers{
		userHandler: users.NewHandler(
			s.router,
			s.store,
			s.cfg.SvixSecret,
		),
		workflow: workflow.NewHandler(
			s.router,
			s.workflowManager,
			s.temporalClient.ScheduleClient(),
			s.middlewareManager,
			s.store,
		),
	}

	s.handlers.userHandler.RegisterRoutes()
	s.handlers.workflow.RegisterRoutes()
}
