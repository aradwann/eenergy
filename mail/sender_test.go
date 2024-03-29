package mail

import (
	"testing"

	"github.com/aradwann/eenergy/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	// skip when the flag is set to prevent the CI from sending emails every time it runs
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig()
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from Eenergy</a></p>
	`
	to := []string{"aradwann@proton.me"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
