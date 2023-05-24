package mail

import (
	"bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	// test won't be executed if the short flag is set
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := "Message form Simple Bank service"

	to := []string{"f.djukic.96@gmail.com"}
	attachedFilesNames := []string{"../wait-for.sh"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachedFilesNames)

	require.NoError(t, err)
}
