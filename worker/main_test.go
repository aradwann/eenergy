package worker

import (
	"testing"

	"github.com/aradwann/eenergy/mail"
	db "github.com/aradwann/eenergy/repository/store"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

func newTestTaskProcessor(t *testing.T, store db.Store, mailer mail.EmailSender) TaskProcessor {

	taskProcessor := NewRedisTaskProcessor(asynq.RedisClientOpt{}, store, mailer)
	require.NotNil(t, taskProcessor)
	return taskProcessor
}
