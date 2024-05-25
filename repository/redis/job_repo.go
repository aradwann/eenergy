package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aradwann/eenergy/worker"
	"github.com/hibiken/asynq"
)

type JobRepository interface {
	EnqueueVerificationEmail(ctx context.Context, username string, email string) error
}

// jobRepository handles operations related to jobs in the job queue.
type jobRepository struct {
	client *asynq.Client
	logger *slog.Logger
}

// NewJobRepository creates a new JobRepository.
func NewJobRepository(redisAddr string, logger *slog.Logger) JobRepository {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &jobRepository{
		client: client,
		logger: logger,
	}
}

// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ JobRepository = (*jobRepository)(nil)

// EnqueueVerificationEmail enqueues a task to send a verification email.
func (r *jobRepository) EnqueueVerificationEmail(ctx context.Context, username string, email string) error {
	payload := worker.EmailVerificationPayload{UserID: username, Email: email}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(worker.TypeEmailVerification, jsonPayload)
	// If I want to change the priority of this task I can change it here
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second), // make room for the DB to commit the transaction before the task is picked up by the worker, otherwise the worker might not find the record
		asynq.Queue(worker.QueueCritical),
	}

	info, err := r.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		r.logger.Error("Failed to enqueue verification email", err)
		return err
	}

	r.logger.LogAttrs(context.Background(),
		slog.LevelInfo,
		"enqueued task",
		slog.String("type", task.Type()),
		slog.String("payload", string(task.Payload())),
		slog.String("queue", info.Queue),
		slog.Int("max_retry", info.MaxRetry),
	)
	return nil
}
