package mail

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"time"

	"github.com/jaytaylor/html2text"
	goMail "gopkg.in/mail.v2"
)

var (
	errUnableConnectSMTP = func(err error) error { return fmt.Errorf("unable to connect to the SMTP server. %w", err) }
)

const (
	TLS      = "TLS"
	StartTLS = "STARTTLS"
)

type SMTPConfig struct {
	ConnectionSecurity                string
	SkipServerCertificateVerification bool
	Hostname                          string
	Server                            string
	Port                              string
	ServerTimeout                     int
	Username                          string
	Password                          string
	EnableSMTPAuth                    bool
	FeedbackName                      string
	FeedbackEmail                     string
	ReplyToAddress                    string
}

type mailData struct {
	mimeTo          []string
	smtpTo          []string
	from            mail.Address
	cc              []string
	replyTo         mail.Address
	subject         string
	htmlBody        string
	txtBody         string
	embeddedFiles   map[string]io.Reader
	attachmentFiles map[string]io.Reader
	mimeHeaders     map[string]string
	messageID       []string
	inReplyTo       []string
	references      []string
}

// smtpClient is implemented by a smtp.Client. See https://golang.org/pkg/net/smtp/#Client.
type smtpClient interface {
	Mail(string) error
	Rcpt(string) error
	Data() (io.WriteCloser, error)
}

func encodeRFC2047Word(s string) string {
	return mime.BEncoding.Encode("utf-8", s)
}

type authChooser struct {
	smtp.Auth
	config *SMTPConfig
}

func (a *authChooser) Start(server *smtp.ServerInfo) (string, []byte, error) {
	smtpAddress := a.config.Server + ":" + a.config.Port
	a.Auth = LoginAuth(a.config.Username, a.config.Password, smtpAddress)
	for _, method := range server.Auth {
		if method == "PLAIN" {
			a.Auth = smtp.PlainAuth("", a.config.Username, a.config.Password, a.config.Server+":"+a.config.Port)
			break
		}
	}
	return a.Auth.Start(server)
}

type loginAuth struct {
	username, password, host string
}

func LoginAuth(username, password, host string) smtp.Auth {
	return &loginAuth{username, password, host}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS {
		return "", nil, fmt.Errorf("unencrypted connection")
	}

	if server.Name != a.host {
		return "", nil, fmt.Errorf("wrong host name")
	}

	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown fromServer")
		}
	}
	return nil, nil
}

func ConnectToSMTPServerAdvanced(config *SMTPConfig) (net.Conn, error) {
	var conn net.Conn
	var err error

	smtpAddress := config.Server + ":" + config.Port
	dialer := &net.Dialer{
		Timeout: time.Duration(config.ServerTimeout) * time.Second,
	}

	if config.ConnectionSecurity == TLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: config.SkipServerCertificateVerification,
			ServerName:         config.Server,
		}

		conn, err = tls.DialWithDialer(dialer, "tcp", smtpAddress, tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to the SMTP server through TLS. %w", err)
		}
		return conn, err
	}
	return dialer.Dial("tcp", smtpAddress)
}

func ConnectToSMTPServer(config *SMTPConfig) (net.Conn, error) {
	return ConnectToSMTPServerAdvanced(config)
}

func NewSMTPClientAdvanced(ctx context.Context, conn net.Conn, config *SMTPConfig) (*smtp.Client, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var c *smtp.Client
	ec := make(chan error)
	go func() {
		var err error
		c, err = smtp.NewClient(conn, config.Server+":"+config.Port)
		if err != nil {
			ec <- err
			return
		}
		cancel()
	}()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil && !errors.Is(err, context.Canceled) {
			return nil, errUnableConnectSMTP(err)
		}
	case err := <-ec:
		return nil, errUnableConnectSMTP(err)
	}

	if config.Hostname != "" {
		err := c.Hello(config.Hostname)
		if err != nil {
			return nil, errUnableConnectSMTP(err)
		}
	}

	if config.ConnectionSecurity == StartTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: config.SkipServerCertificateVerification,
			ServerName:         config.Server,
		}
		if err := c.StartTLS(tlsConfig); err != nil {
			return nil, err
		}
	}

	if config.EnableSMTPAuth {
		if err := c.Auth(&authChooser{config: config}); err != nil {
			return nil, fmt.Errorf("authentication failed. %w", err)
		}
	}
	return c, nil
}

func NewSMTPClient(ctx context.Context, conn net.Conn, config *SMTPConfig) (*smtp.Client, error) {
	return NewSMTPClientAdvanced(
		ctx,
		conn,
		config,
	)
}

func TestConnection(config *SMTPConfig) error {
	conn, err := ConnectToSMTPServer(config)
	if err != nil {
		return errUnableConnectSMTP(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	sec := config.ServerTimeout

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(sec)*time.Second)
	defer cancel()

	c, err := NewSMTPClient(ctx, conn, config)
	if err != nil {
		return errUnableConnectSMTP(err)
	}
	_ = c.Close()
	_ = c.Quit()

	return nil
}

func SendMailWithAttachmentFilesUsingConfig(ctx context.Context, to []string, subject, htmlBody, txtBody string, attachmentFiles map[string]io.Reader, config *SMTPConfig, messageID []string, inReplyTo []string, references []string, ccMail []string) error {
	fromMail := mail.Address{Name: config.FeedbackName, Address: config.FeedbackEmail}
	replyTo := mail.Address{Name: config.FeedbackName, Address: config.ReplyToAddress}

	data := mailData{
		mimeTo:          to,
		smtpTo:          to,
		from:            fromMail,
		cc:              ccMail,
		replyTo:         replyTo,
		subject:         subject,
		htmlBody:        htmlBody,
		txtBody:         txtBody,
		attachmentFiles: attachmentFiles,
		messageID:       messageID,
		inReplyTo:       inReplyTo,
		references:      references,
	}

	return sendMailUsingConfigAdvanced(ctx, data, config)
}

func SendMailUsingConfig(ctx context.Context, to []string, subject, htmlBody, txtBody string, config *SMTPConfig, messageID, inReplyTo, references, ccMail []string) error {
	return SendMailWithAttachmentFilesUsingConfig(ctx, to, subject, htmlBody, txtBody, nil, config, messageID, inReplyTo, references, ccMail)
}

// allows for sending an email with differing MIME/SMTP recipients
func sendMailUsingConfigAdvanced(ctx context.Context, mail mailData, config *SMTPConfig) error {
	if config.Server == "" {
		return nil
	}

	conn, err := ConnectToSMTPServer(config)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	sec := config.ServerTimeout

	ctx, cancel := context.WithTimeout(ctx, time.Duration(sec)*time.Second)
	defer cancel()

	c, err := NewSMTPClient(ctx, conn, config)
	if err != nil {
		return err
	}
	defer func() {
		_ = c.Quit()
	}()
	defer func() {
		_ = c.Close()
	}()

	return SendMail(c, mail, time.Now())
}

func SendMail(c smtpClient, mail mailData, date time.Time) error {
	htmlMessage, txtBody, headers := buildMailMsg(mail)

	m := goMail.NewMessage(goMail.SetCharset("UTF-8"))
	m.SetHeaders(headers)
	m.SetDateHeader("Date", date)
	m.SetBody("text/plain", txtBody)
	if len(htmlMessage) > 0 {
		m.AddAlternative("text/html", htmlMessage)
	}

	for name, reader := range mail.attachmentFiles {
		m.AttachReader(name, reader)
	}

	if err := c.Mail(mail.from.Address); err != nil {
		return fmt.Errorf("failed to set the from address. %w", err)
	}

	for _, s := range mail.smtpTo {
		if err := c.Rcpt(s); err != nil {
			return fmt.Errorf("failed to set the to address. %w", err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to add email message data. %w", err)
	}

	_, err = m.WriteTo(w)
	if err != nil {
		return fmt.Errorf("failed to write the email message. %w", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection to the SMTP server. %w", err)
	}

	return nil
}

func buildMailMsg(mail mailData) (string, string, map[string][]string) {
	htmlMessage := mail.htmlBody
	txtBody := mail.txtBody
	if len(txtBody) == 0 {
		txtBody, _ = html2text.FromString(mail.htmlBody)
	}

	headers := map[string][]string{
		"From":                      {mail.from.String()},
		"To":                        mail.mimeTo,
		"Subject":                   {encodeRFC2047Word(mail.subject)},
		"Content-Transfer-Encoding": {"8bit"},
		"Auto-Submitted":            {"auto-generated"},
		"Precedence":                {"bulk"},
	}

	if mail.replyTo.Address != "" {
		headers["Reply-To"] = []string{mail.replyTo.String()}
	}

	if len(mail.cc) > 0 {
		headers["CC"] = mail.cc
	}

	if len(mail.messageID) > 0 {
		headers["Message-ID"] = mail.messageID
	}

	if len(mail.inReplyTo) > 0 {
		headers["In-Reply-To"] = mail.inReplyTo
	}

	if len(mail.references) > 0 {
		headers["References"] = mail.references
	}

	for k, v := range mail.mimeHeaders {
		headers[k] = []string{encodeRFC2047Word(v)}
	}
	return htmlMessage, txtBody, headers
}
