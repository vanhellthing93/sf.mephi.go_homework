package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/vanhellthing93/sf.mephi.go_homework/config"
)

type SMTPService struct {
	dialer *mail.Dialer
	from   string
}

func NewSMTPService() *SMTPService {
	// Загружаем конфигурацию SMTP
	smtpConfig := config.LoadSMTPConfig()

	// Создаем dialer
	dialer := mail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password)
	dialer.TLSConfig = &tls.Config{
		ServerName:         smtpConfig.Host,
		InsecureSkipVerify: false,
	}

	return &SMTPService{
		dialer: dialer,
		from:   smtpConfig.From,
	}
}

func (s *SMTPService) SendEmail(to, subject, body string) error {
	// Создаем сообщение
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Отправляем сообщение
	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("SMTP error: %v", err)
		return fmt.Errorf("email sending failed")
	}

	log.Printf("Email sent to %s", to)
	return nil
}

func (s *SMTPService) SendRegistrationNotification(userEmail string) error {
	// Создаем тело письма
	content := fmt.Sprintf(`
		<h1>Добро пожаловать в наш банковский сервис!</h1>
		<p>Ваш аккаунт успешно создан.</p>
		<p>Дата: %s</p>
		<small>Это автоматическое уведомление</small>
	`, time.Now().Format("02.01.2006 15:04:05"))

	// Отправляем письмо
	return s.SendEmail(userEmail, "Регистрация в банковском сервисе", content)
}

func (s *SMTPService) SendOverduePaymentNotification(userEmail string, amount float64) error {
	// Создаем тело письма
	content := fmt.Sprintf(`
		<h1>У вас есть просроченный платеж!</h1>
		<p>Сумма: <strong>%.2f RUB</strong></p>
		<p>Пожалуйста, погасите задолженность как можно скорее.</p>
		<p>Дата: %s</p>
		<small>Это автоматическое уведомление</small>
	`, amount, time.Now().Format("02.01.2006 15:04:05"))

	// Отправляем письмо
	return s.SendEmail(userEmail, "Просроченный платеж", content)
}

func (s *SMTPService) SendCreditNotification(userEmail string, amount, interestRate float64, term int) error {
	// Создаем тело письма
	content := fmt.Sprintf(`
		<h1>Кредит успешно оформлен!</h1>
		<p>Сумма: <strong>%.2f RUB</strong></p>
		<p>Процентная ставка: <strong>%.2f%%</strong></p>
		<p>Срок: <strong>%d месяцев</strong></p>
		<p>Дата: %s</p>
		<small>Это автоматическое уведомление</small>
	`, amount, interestRate, term, time.Now().Format("02.01.2006 15:04:05"))

	// Отправляем письмо
	return s.SendEmail(userEmail, "Кредит успешно оформлен", content)
}