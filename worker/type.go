package worker

const (
	TypeEmailVerification = "email:verification"
)

type EmailVerificationPayload struct {
	UserID string
	Email  string
}
