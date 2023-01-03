package worker

import (
	db "backend_masterclass/db/sqlc"
	"context"

	"github.com/hibiken/asynq"
)

const (
	CriticalQueue       = "critical"
	DefaultQueue        = "default"
	CriticalQueryWeight = 2
)

// This file is supposed to take tasks and assign them to the Redis worker
type TaskProcessor interface {
	//It is important, very important to register the task with our processor
	//Before starting the server
	Start() error
	ProcessTask(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {

	sever := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				CriticalQueue: 10,
				DefaultQueue:  5,
			},
		},
	)

	return &RedisTaskProcessor{
		server: sever,
		store:  store,
	}

}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTask)

	return processor.server.Start(mux)
}
