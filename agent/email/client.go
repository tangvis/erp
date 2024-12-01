package email

import (
	"context"
	"io"
	"os"

	mail "github.com/tangvis/erp/pkg/emailhelper"
)

type client struct {
	config *mail.SMTPConfig
}

// NewDefaultClient
// 默认通用发件人
func NewDefaultClient() Client {
	return NewClient(&mail.SMTPConfig{
		Server:        "smtp.google.com",
		Port:          "587",
		ServerTimeout: 10,
		FeedbackName:  "SPX System Email",
		FeedbackEmail: "stats.spx@google.com",
	})
}

func NewClient(config *mail.SMTPConfig) Client {
	return &client{
		config: config,
	}
}

func (c *client) SendMailCommon(ctx context.Context, to, ccMail, messageID, inReplyTo, references []string, subject, htmlBody, txtBody string) error {
	return mail.SendMailUsingConfig(ctx, to, subject, htmlBody, txtBody, c.config, messageID, inReplyTo, references, ccMail)
}

func (c *client) SendMailWithCC(ctx context.Context, to, ccMail []string, subject, htmlBody, txtBody string) error {
	return c.SendMailCommon(ctx, to, ccMail, []string{}, []string{}, []string{}, subject, htmlBody, txtBody)
}

func (c *client) SendHTMLMail(ctx context.Context, to []string, subject, htmlBody string) error {
	return c.SendMailWithCC(ctx, to, []string{}, subject, htmlBody, "")
}

func (c *client) SendTxtMail(ctx context.Context, to []string, subject, txtBody string) error {
	return c.SendMailWithCC(ctx, to, []string{}, subject, "", txtBody)
}

func (c *client) SendMailWithAttachmentFiles(ctx context.Context, to []string, subject, htmlBody, txtBody string, attachmentFile map[string]io.Reader) error {
	return c.SendMailWithAttachmentFilesCommon(ctx, to, []string{}, []string{}, []string{}, subject, htmlBody, txtBody, attachmentFile)
}

func (c *client) SendMailWithFilepath(ctx context.Context, to []string, subject, htmlBody, txtBody string, filepath string) error {
	fileReader, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = fileReader.Close()
	}()
	return c.SendMailWithAttachmentFile(ctx, to, subject, htmlBody, txtBody, fileReader.Name(), fileReader)
}

func (c *client) SendMailWithAttachmentFile(ctx context.Context, to []string, subject, htmlBody, txtBody string, filename string, fileReader io.Reader) error {
	if len(filename) == 0 {
		filename = "default_filename"
	}
	return c.SendMailWithAttachmentFilesCommon(ctx, to, []string{}, []string{}, []string{}, subject, htmlBody, txtBody, map[string]io.Reader{filename: fileReader})
}

func (c *client) SendMailWithAttachmentFilesCommon(ctx context.Context, to, messageID, inReplyTo, references []string, subject, htmlBody, txtBody string, attachmentFile map[string]io.Reader) error {
	return mail.SendMailWithAttachmentFilesUsingConfig(ctx, to, subject, htmlBody, txtBody, attachmentFile, c.config, messageID, inReplyTo, references, []string{})
}
