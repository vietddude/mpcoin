package mail

import (
	"fmt"
	"mpc/internal/infrastructure/config"

	gomail "github.com/wneessen/go-mail"
)

type Client struct {
	client *gomail.Client
	from   string
}

func NewClient(cfg *config.MailConfig) (*Client, error) {
	client, err := gomail.NewClient(cfg.SMTPHost,
		gomail.WithPort(cfg.SMTPPort),
		gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
		gomail.WithUsername(cfg.SMTPUsername),
		gomail.WithPassword(cfg.SMTPPassword),
		gomail.WithTLSPolicy(gomail.TLSMandatory), // Adjust based on your SMTP server requirements
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	return &Client{
		client: client,
		from:   cfg.FromEmail,
	}, nil
}

func (c *Client) SendMail(to, subject, body string) error {
	m := gomail.NewMsg()
	m.From(c.from)
	m.To(to)
	m.Subject(subject)
	m.SetBodyString(gomail.TypeTextHTML, body)

	if err := c.client.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
