package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type SMTPSender struct {
	host string
	port string
	user string
	pass string
	from string
	auth smtp.Auth
}

func NewSMTPSender() *SMTPSender {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if host == "" {
		host = "smtp.gmail.com"
	}
	if port == "" {
		port = "587"
	}

	return &SMTPSender{
		host: host,
		port: port,
		user: user,
		pass: pass,
		from: from,
		auth: smtp.PlainAuth("", user, pass, host),
	}
}

func (s *SMTPSender) send(to, subject, body string) error {
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.from, to, subject, body,
	)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, s.auth, s.user, []string{to}, []byte(msg))
}

func (s *SMTPSender) SendWelcome(toEmail, toName string) error {
	body := fmt.Sprintf(`
	<html><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px">
	  <div style="background:#534AB7;color:#fff;padding:24px;border-radius:8px 8px 0 0;text-align:center">
	    <h1 style="margin:0">Welcome to Cinema!</h1>
	  </div>
	  <div style="background:#f9f9f9;padding:24px;border-radius:0 0 8px 8px">
	    <p>Hi <strong>%s</strong>,</p>
	    <p>Your account has been created successfully. You can now browse movies, book tickets and manage your reservations.</p>
	    <p style="color:#888;font-size:13px">If you didn't create this account, please ignore this email.</p>
	  </div>
	</body></html>`, toName)
	return s.send(toEmail, "Welcome to Cinema System!", body)
}

func (s *SMTPSender) SendBookingConfirmed(toEmail, toName, movieTitle, showtime string, seats []string, amount float64) error {
	seatList := strings.Join(seats, ", ")
	body := fmt.Sprintf(`
	<html><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px">
	  <div style="background:#534AB7;color:#fff;padding:24px;border-radius:8px 8px 0 0;text-align:center">
	    <h1 style="margin:0">Booking Confirmed!</h1>
	  </div>
	  <div style="background:#f9f9f9;padding:24px;border-radius:0 0 8px 8px">
	    <p>Hi <strong>%s</strong>, your booking is confirmed.</p>
	    <table style="width:100%%;border-collapse:collapse;margin:16px 0">
	      <tr style="background:#eee"><td style="padding:8px;font-weight:bold">Movie</td><td style="padding:8px">%s</td></tr>
	      <tr><td style="padding:8px;font-weight:bold">Showtime</td><td style="padding:8px">%s</td></tr>
	      <tr style="background:#eee"><td style="padding:8px;font-weight:bold">Seats</td><td style="padding:8px">%s</td></tr>
	      <tr><td style="padding:8px;font-weight:bold">Total</td><td style="padding:8px;color:#534AB7;font-size:18px"><strong>$%.2f</strong></td></tr>
	    </table>
	    <p style="color:#888;font-size:13px">Please arrive 15 minutes before the showtime. Enjoy the movie!</p>
	  </div>
	</body></html>`, toName, movieTitle, showtime, seatList, amount)
	return s.send(toEmail, fmt.Sprintf("Booking Confirmed — %s", movieTitle), body)
}
