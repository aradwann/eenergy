package worker

import (
	"testing"

	db "github.com/aradwann/eenergy/db/store"
	"github.com/aradwann/eenergy/mail"
	"github.com/hibiken/asynq"
)

func newTestTaskProcessor(t *testing.T, store db.Store, mailer mail.EmailSender) TaskProcessor {

	taskProcessor := NewRedisTaskProcessor(asynq.RedisClientOpt{}, store, mailer)

	return taskProcessor
}
