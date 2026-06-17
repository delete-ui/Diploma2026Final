package email

import (
	"GolangBackendDiploma26/internal/models"
	"fmt"
	"go.uber.org/zap"

	"gopkg.in/gomail.v2"
)

type Sender interface {
	SendVerificationCode(to, username, code string) error
	SendPasswordResetCode(to, username, code string) error
	SendReceipt(to, username, body string) error
}

type SMTPSender struct {
	config models.SMTPConfig
	dialer *gomail.Dialer
}

func NewSMTPSender(cfg models.SMTPConfig) (*SMTPSender, error) {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return &SMTPSender{
		config: cfg,
		dialer: dialer,
	}, nil
}

func (s *SMTPSender) SendVerificationCode(to, username, code string) error {
	subject := "Подтверждение регистрации в Battery Shop"
	body := fmt.Sprintf(`
		<h2>Здравствуйте, %s!</h2>
		<p>Спасибо за регистрацию в магазине автомобильных аккумуляторов.</p>
		<p>Ваш код подтверждения: <strong>%s</strong></p>
		<p>Код действителен 1 час.</p>
	`, username, code)
	return s.sendEmail(to, subject, body)
}

func (s *SMTPSender) SendPasswordResetCode(to, username, code string) error {
	subject := "Сброс пароля в Battery Shop"
	body := fmt.Sprintf(`
		<h2>Здравствуйте, %s!</h2>
		<p>Вы запросили сброс пароля.</p>
		<p>Ваш код для сброса: <strong>%s</strong></p>
		<p>Код действителен 1 час.</p>
	`, username, code)
	return s.sendEmail(to, subject, body)
}

func (s *SMTPSender) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	from := s.config.From
	if s.config.FromName != "" && s.config.FromName != s.config.From {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.From)
	}
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}

type LogSender struct {
	logger *zap.Logger
}

func NewLogSender(logger *zap.Logger) *LogSender {
	return &LogSender{logger: logger}
}

func (l *LogSender) SendVerificationCode(to, username, code string) error {
	l.logger.Info("verification code (email stub)",
		zap.String("to", to),
		zap.String("username", username),
		zap.String("code", code),
	)
	return nil
}

func (s *SMTPSender) SendReceipt(to, username, body string) error {
	subject := fmt.Sprintf("Заказ в Battery Shop")
	return s.sendEmail(to, subject, body)
}

func (l *LogSender) SendReceipt(to, username, body string) error {
	l.logger.Info("receipt (email stub)", zap.String("to", to))
	return nil
}

func (l *LogSender) SendPasswordResetCode(to, username, code string) error {
	l.logger.Info("password reset code (email stub)",
		zap.String("to", to),
		zap.String("username", username),
		zap.String("code", code),
	)
	return nil
}
