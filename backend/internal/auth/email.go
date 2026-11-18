package auth

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	host string
	port int
	from string
}

func NewMailer(host string, port int, from string) *Mailer {
	return &Mailer{
		host: host,
		port: port,
		from: from,
	}
}

func (m *Mailer) SendPasswordConfirmEmail(toEmail, token string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", "Подтверждение смены пароля в Astro-Planner")

	confirmURL := fmt.Sprintf("http://localhost:8080/api/profile/password/confirm?token=%s", token)
	body := fmt.Sprintf(`
		<h3>Запрос на изменение пароля</h3>
		<p>Для подтверждения нового пароля перейдите по ссылке (действует 15 минут):</p>
		<a href="%s">Подтвердить новый пароль</a>
		<p>Если это были не вы — просто проигнорируйте это письмо.</p>
	`, confirmURL)

	msg.SetBody("text/html", body)

	d := gomail.NewDialer(m.host, m.port, "", "")
	return d.DialAndSend(msg)
}
