package worker

import (
	"time"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"stock-agent.io/configs"
)

type Worker struct {
	worker         worker.Worker
	temporalClient client.Client
	taskQueue      string
}

func New(cfg *configs.AppConfig, temporalClient client.Client, taskQueue string) *Worker {
	workerOptions := worker.Options{
		MaxConcurrentActivityExecutionSize:      10,
		MaxConcurrentWorkflowTaskExecutionSize:  10,
		MaxConcurrentLocalActivityExecutionSize: 10,
		EnableLoggingInReplay:                   true,
		Identity:                                cfg.AppName,
		DeadlockDetectionTimeout:                time.Minute,
		MaxHeartbeatThrottleInterval:            time.Second * 60,
		DefaultHeartbeatThrottleInterval:        time.Second * 30,
	}

	w := worker.New(temporalClient, taskQueue, workerOptions)

	return &Worker{
		worker:         w,
		temporalClient: temporalClient,
		taskQueue:      taskQueue,
	}
}

func (w *Worker) RegisterWorkflow(workflow interface{}) {
	w.worker.RegisterWorkflow(workflow)
	log.Info().Msgf("Registered workflow: %T", workflow)
}

func (w *Worker) RegisterActivity(activity interface{}) {
	w.worker.RegisterActivity(activity)
	log.Info().Msgf("Registered activity: %T", activity)
}

func (w *Worker) Start() error {
	log.Info().Str("task_queue", w.taskQueue).Msg("Starting Temporal worker")
	return w.worker.Run(worker.InterruptCh())
}

func (w *Worker) Stop() {
	log.Info().Msg("Stopping Temporal worker")
	w.worker.Stop()
}
