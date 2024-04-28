package service

import (
	"context"

	"github.com/tangvis/erp/app/system/email/define"
)

type APP interface {
	Send(ctx context.Context, mail define.MailInfo) error
}
