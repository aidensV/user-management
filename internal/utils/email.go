package utils

import (
	"fmt"
	"net/smtp"
)

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

var emailConfig *EmailConfig

// InitEmailConfig initializes email configuration
func InitEmailConfig(host, port, username, password, from string) {
	emailConfig = &EmailConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

// SendResetPasswordEmail sends password reset email
func SendResetPasswordEmail(to, token, resetURL string) error {
	if emailConfig == nil {
		// Skip if email not configured
		return nil
	}

	link := fmt.Sprintf("%s?token=%s", resetURL, token)
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
        <html>
        <body>
            <h2>Password Reset Request</h2>
            <p>You requested to reset your password. Click the link below to reset it:</p>
            <p><a href="%s">%s</a></p>
            <p>This link will expire in 1 hour.</p>
            <p>If you did not request this, please ignore this email.</p>
            <br>
            <p>Best regards,</p>
            <p>User Management System</p>
        </body>
        </html>
    `, link, link)

	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body)

	addr := fmt.Sprintf("%s:%s", emailConfig.Host, emailConfig.Port)
	return smtp.SendMail(addr, auth, emailConfig.From, []string{to}, msg)
}
