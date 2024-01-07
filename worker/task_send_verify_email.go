package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verfiy_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	slog.LogAttrs(context.Background(),
		slog.LevelInfo,
		"enqueued task",
		slog.String("type", task.Type()),
		slog.String("payload", string(task.Payload())),
		slog.String("queue", info.Queue),
		slog.Int("max_retry", info.MaxRetry),
	)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		// if the user is not found try again later, as the creation might be not commited yet
		// if errors.Is(err, sql.ErrNoRows) {
		// 	return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)

		// }
		return fmt.Errorf("failed to get user: %w", asynq.SkipRetry)
	}
	// TODO: send email to user
	slog.LogAttrs(context.Background(),
		slog.LevelInfo,
		"processed task",
		slog.String("type", task.Type()),
		slog.String("payload", string(task.Payload())),
		slog.String("email", user.Email),
	)
	return nil
}
