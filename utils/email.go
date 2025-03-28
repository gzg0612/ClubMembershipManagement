package utils

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// EmailClient 邮件客户端
type EmailClient struct {
	config *EmailConfig
}

// NewEmailClient 创建邮件客户端
func NewEmailClient(config *EmailConfig) *EmailClient {
	return &EmailClient{
		config: config,
	}
}

// Send 发送普通文本邮件
func (c *EmailClient) Send(to, subject, content string) error {
	e := email.NewEmail()
	e.From = c.config.From
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(content)

	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)

	return e.Send(addr, auth)
}

// SendHTML 发送HTML邮件
func (c *EmailClient) SendHTML(to, subject, htmlContent string) error {
	e := email.NewEmail()
	e.From = c.config.From
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(htmlContent)

	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)

	return e.Send(addr, auth)
}
