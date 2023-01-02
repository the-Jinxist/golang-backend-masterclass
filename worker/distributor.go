package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	asynq "github.com/hibiken/asynq"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	byteArray, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error while marshalling payload: %s", err.Error())
	}

	task := asynq.NewTask(TaskSendVerifyEmail, byteArray, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %d: ", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_entry", info.MaxRetry).
		Msg("Enqueued task!")

	return nil

}

func (processor *RedisTaskProcessor) ProcessTask(ctx context.Context, task *asynq.Task) error {

	var payload PayloadSendVerifyEmail

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		//We call `asynq.SkipRetry` to tell Redis there is no need to retry the task
		return fmt.Errorf("failed to unmarshal payload: %d: ", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			//Also, in this case there is no need to retry
			return fmt.Errorf("this user doesn't exist: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get user: %s", err.Error())
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task!")

	return nil
}
