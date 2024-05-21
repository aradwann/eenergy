package worker

// const (
// 	QueueCritical = "critical"
// 	QueueDefault  = "default"
// 	QueueLow      = "low"
// )

// type TaskProcessor interface {
// 	Start() error
// 	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
// }

// type RedisTaskProcessor struct {
// 	server *asynq.Server
// 	store  db.Store
// 	mailer mail.EmailSender
// }

// func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
// 	server := asynq.NewServer(redisOpt, asynq.Config{
// 		Queues: map[string]int{
// 			QueueCritical: 6,
// 			QueueDefault:  3,
// 			QueueLow:      1,
// 		},
// 		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
// 			slog.LogAttrs(ctx,
// 				slog.LevelError,
// 				"process task failed",
// 				slog.String("err", err.Error()),
// 				slog.String("type", task.Type()),
// 				slog.String("payload", string(task.Payload())))
// 		}),
// 		Logger: NewLogger(),
// 	})

// 	return &RedisTaskProcessor{
// 		server: server,
// 		store:  store,
// 		mailer: mailer,
// 	}
// }

// func (processor *RedisTaskProcessor) Start() error {
// 	mux := asynq.NewServeMux()
// 	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
// 	return processor.server.Start(mux)
// }
