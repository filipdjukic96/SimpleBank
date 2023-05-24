package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachedFilesNames []string) error
}

type GmailSender struct {
	name              string //sender of the email
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachedFilesNames []string) error {
	email := email.NewEmail()
	email.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	email.Subject = subject
	email.HTML = []byte(content)
	email.To = to
	email.Cc = cc
	email.Bcc = bcc

	for _, f := range attachedFilesNames {
		_, err := email.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	// authentication
	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return email.Send(smtpServerAddress, smtpAuth)
}
