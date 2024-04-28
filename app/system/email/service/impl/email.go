package impl

import (
	"context"
	"html/template"
	"strings"
	"time"

	"github.com/tangvis/erp/agent/email"
	"github.com/tangvis/erp/agent/templates"
	"github.com/tangvis/erp/app/system/email/define"
	"github.com/tangvis/erp/app/system/email/repository"
	"github.com/tangvis/erp/app/system/email/service"
)

type Email struct {
	mailClient email.Client
	repo       repository.Repo
	container  *templates.Container
}

func NewEmailAPP(
	mailClient email.Client,
	repo repository.Repo,
	container *templates.Container,
) service.APP {
	return &Email{
		mailClient: mailClient,
		repo:       repo,
		container:  container,
	}
}

func (e Email) Send(ctx context.Context, mail define.MailInfo) error {
	if err := mail.Validate(); err != nil {
		return err
	}
	var (
		err    error
		start  = time.Now()
		record = repository.EmailRecordTab{
			Operator:     mail.Operator,
			Sender:       "system",
			Receivers:    strings.Join(mail.To, ","),
			TemplateName: mail.Template,
			Subject:      mail.Subject,
			SendStatus:   define.Init,
		}
	)
	defer func() {
		record.ExecutionTime = time.Since(start).Milliseconds()
		record.SendStatus = define.Send
		if err != nil {
			record.Result = err.Error()
			record.SendStatus = define.Failed
		}
		err = e.repo.Save(ctx, record)
	}()
	mailContent, err := e.container.RenderToString(mail.Template, templates.Data{
		Props: mail.Content,
		HTML:  make(map[string]template.HTML),
	})
	if err != nil {
		return err
	}
	record.Content = mailContent
	err = e.mailClient.SendHTMLMail(ctx, mail.To, mail.Subject, mailContent)

	return err
}
