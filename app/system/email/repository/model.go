package repository

import (
	"github.com/tangvis/erp/agent/mysql"
	"github.com/tangvis/erp/app/system/email/define"
)

type EmailRecordTab struct {
	Operator      string
	Sender        string
	Receivers     string // seperated by ,
	TemplateName  string
	Attachment    string
	Subject       string
	Content       string
	ExecutionTime int64 // ms
	SendStatus    define.Status
	Result        string
	mysql.BaseModel
}
