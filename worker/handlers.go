package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aradwann/eenergy/mail"
	emailDB "github.com/aradwann/eenergy/repository/postgres/email"
	userDB "github.com/aradwann/eenergy/repository/postgres/user"
	"github.com/aradwann/eenergy/util"
	"github.com/hibiken/asynq"
)

type EmailHandler interface {
	HandleEmailVerificationTask(ctx context.Context, task *asynq.Task) error
}

type redisEmailHandler struct {
	userRepo  userDB.UserRepository
	emailRepo emailDB.EmailRepository
	mailer    mail.EmailSender
	logger    *slog.Logger
}

var _ EmailHandler = (*redisEmailHandler)(nil)

func NewEmailHandler(userRepo userDB.UserRepository, emailRepo emailDB.EmailRepository, mailer mail.EmailSender, logger *slog.Logger) EmailHandler {
	return &redisEmailHandler{
		userRepo:  userRepo,
		emailRepo: emailRepo,
		mailer:    mailer,
		logger:    logger,
	}
}

func (h *redisEmailHandler) HandleEmailVerificationTask(ctx context.Context, task *asynq.Task) error {
	var p EmailVerificationPayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		h.logger.Error("Failed to unmarshal payload", err)
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	// Send email verification logic here
	h.logger.Info("Sending email verification", "userID", p.UserID, "email", p.Email)
	user, err := h.userRepo.GetUser(ctx, p.UserID)
	if err != nil {
		// TODO: handle not found SQL error code
		return fmt.Errorf("gailed to get user: %w", err)
	}
	createVerifyEmailParams := emailDB.CreateVerifyEmailParams{
		Username:   p.UserID,
		Email:      p.Email,
		SecretCode: util.RandomString(32),
	}
	verfiyEmail, err := h.emailRepo.CreateVerifyEmail(ctx, createVerifyEmailParams)
	if err != nil {
		return fmt.Errorf("failed to create verify email instance: %w", err)
	}
	subject := "Welcome to Eennergy"
	verifyURL := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", verfiyEmail.ID, verfiyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s, <br/>
	Thank you for being a member in Eenergy community!</br>
	Pleas click on <a href="%s">click here</a> to verify your email`, user.FullName, verifyURL)
	to := []string{user.Email}
	err = h.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}
	slog.LogAttrs(context.Background(),
		slog.LevelInfo,
		"processed task",
		slog.String("type", task.Type()),
		slog.String("payload", string(task.Payload())),
		slog.String("email", user.Email),
	)
	return nil
}
