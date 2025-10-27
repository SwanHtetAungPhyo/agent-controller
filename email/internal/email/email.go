package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"
)

type Config struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	UseTLS       bool
	UseSSL       bool
}

type EmailService struct {
	config Config
	auth   smtp.Auth
}

type Email struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	HTMLBody    string
	Attachments []Attachment
	ReplyTo     string
}

type Attachment struct {
	Filename string
	Content  []byte
	MimeType string
}

func NewEmailService(cfg Config) *EmailService {
	auth := smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)

	return &EmailService{
		config: cfg,
		auth:   auth,
	}
}

func (es *EmailService) SendEmail(email *Email) error {
	if len(email.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	message, err := es.buildMessage(email)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", es.config.SMTPHost, es.config.SMTPPort)

	recipients := append(email.To, email.Cc...)
	recipients = append(recipients, email.Bcc...)

	if es.config.UseSSL {
		return es.sendWithSSL(addr, recipients, message)
	} else if es.config.UseTLS {
		return es.sendWithTLS(addr, recipients, message)
	} else {
		return smtp.SendMail(addr, es.auth, es.config.FromEmail, recipients, message)
	}
}

func (es *EmailService) sendWithTLS(addr string, recipients []string, message []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	tlsConfig := &tls.Config{
		ServerName: es.config.SMTPHost,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	if err = client.Auth(es.auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err = client.Mail(es.config.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range recipients {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = writer.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return client.Quit()
}

func (es *EmailService) sendWithSSL(addr string, recipients []string, message []byte) error {
	tlsConfig := &tls.Config{
		ServerName: es.config.SMTPHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to establish SSL connection: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, es.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err = client.Auth(es.auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err = client.Mail(es.config.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range recipients {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = writer.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return client.Quit()
}

func (es *EmailService) buildMessage(email *Email) ([]byte, error) {
	var buf bytes.Buffer

	from := fmt.Sprintf("%s <%s>", es.config.FromName, es.config.FromEmail)
	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))

	if len(email.Cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(email.Cc, ", ")))
	}

	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))

	if email.ReplyTo != "" {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", email.ReplyTo))
	}

	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	buf.WriteString("MIME-Version: 1.0\r\n")

	if email.HTMLBody != "" || len(email.Attachments) > 0 {
		boundary := "----=_Part_0_" + fmt.Sprint(time.Now().Unix())
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
		buf.WriteString("\r\n")

		if email.HTMLBody != "" {
			buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
			buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
			buf.WriteString("\r\n")
			buf.WriteString(email.HTMLBody)
			buf.WriteString("\r\n")
		} else if email.Body != "" {
			buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
			buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
			buf.WriteString("\r\n")
			buf.WriteString(email.Body)
			buf.WriteString("\r\n")
		}

		for _, attachment := range email.Attachments {
			buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", attachment.MimeType, attachment.Filename))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", attachment.Filename))
			buf.WriteString("\r\n")
			buf.Write(attachment.Content)
			buf.WriteString("\r\n")
		}

		buf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(email.Body)
	}

	return buf.Bytes(), nil
}

func (es *EmailService) SendSimpleEmail(to []string, subject, body string) error {
	email := &Email{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	return es.SendEmail(email)
}

func (es *EmailService) SendHTMLEmail(to []string, subject, htmlBody string) error {
	email := &Email{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlBody,
	}
	return es.SendEmail(email)
}

func (es *EmailService) SendTemplateEmail(to []string, subject, templateStr string, data interface{}) error {
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	email := &Email{
		To:       to,
		Subject:  subject,
		HTMLBody: buf.String(),
	}

	return es.SendEmail(email)
}

func (es *EmailService) SendBulkEmail(subject, body string, recipients []string) error {
	email := &Email{
		To:      []string{es.config.FromEmail}, // Send to self
		Bcc:     recipients,                    // BCC all recipients
		Subject: subject,
		Body:    body,
	}
	return es.SendEmail(email)
}

func (es *EmailService) Validate() error {
	if es.config.SMTPHost == "" {
		return fmt.Errorf("SMTP host is required")
	}
	if es.config.SMTPPort == 0 {
		return fmt.Errorf("SMTP port is required")
	}
	if es.config.FromEmail == "" {
		return fmt.Errorf("from email is required")
	}
	return nil
}
