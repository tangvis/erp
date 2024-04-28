package email

import (
	"context"
	"io"
)

type Client interface {
	SendTxtMail(ctx context.Context, to []string, subject, txtBody string) error
	SendHTMLMail(ctx context.Context, to []string, subject, htmlBody string) error
	SendMailWithAttachmentFile(ctx context.Context, to []string, subject, htmlBody, txtBody string, filename string, fileReader io.Reader) error
	SendMailWithCC(ctx context.Context, to, ccMail []string, subject, htmlBody, txtBody string) error

	SendMailCommon(ctx context.Context, to, ccMail, messageID, inReplyTo, references []string, subject, htmlBody, txtBody string) error
	SendMailWithAttachmentFilesCommon(ctx context.Context, to, messageID, inReplyTo, references []string, subject, htmlBody, txtBody string, embeddedFiles map[string]io.Reader) error

	SendMailWithFilepath(ctx context.Context, to []string, subject, htmlBody, txtBody string, filepath string) error
	SendMailWithAttachmentFiles(ctx context.Context, to []string, subject, htmlBody, txtBody string, embeddedFiles map[string]io.Reader) error
}
