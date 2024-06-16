package worker

import (
	"context"
	"log/slog"
	"os"

	"github.com/aradwann/eenergy/mail"
	emailDB "github.com/aradwann/eenergy/repository/postgres/email"
	userDB "github.com/aradwann/eenergy/repository/postgres/user"
	"github.com/aradwann/eenergy/util"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	start() error
}

type RedisTaskProcessor struct {
	server       *asynq.Server
	emailhandler EmailHandler
	logger       *slog.Logger
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, emailhandler EmailHandler, logger *slog.Logger) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 6,
			QueueDefault:  3,
			QueueLow:      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			logger.LogAttrs(ctx,
				slog.LevelError,
				"process task failed",
				slog.String("err", err.Error()),
				slog.String("type", task.Type()),
				slog.String("payload", string(task.Payload())))
		}),
		Logger: NewLogger(logger),
	})

	return &RedisTaskProcessor{
		server:       server,
		emailhandler: emailhandler,
		logger:       logger,
	}
}

func (processor *RedisTaskProcessor) start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailVerification, processor.emailhandler.HandleEmailVerificationTask)
	return processor.server.Start(mux)
}

// StartTaskProcessor runs the task processor.
func StartTaskProcessor(config util.Config, userRepo userDB.UserRepository, emailRepo emailDB.EmailRepository, logger *slog.Logger) {
	redisOpts := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	emailHandler := NewEmailHandler(userRepo, emailRepo, mailer, logger)
	taskProcessor := NewRedisTaskProcessor(redisOpts, emailHandler, logger)
	logger.Info("start task processor")
	err := taskProcessor.start()
	if err != nil {
		logger.Error("failed to start redis Task processor (workers) ", err)
		os.Exit(1)
	}
}
