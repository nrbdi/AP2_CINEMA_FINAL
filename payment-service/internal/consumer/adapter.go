package consumer

type emailAdapter struct {
	sender interface {
		SendPaymentReceipt(toEmail, toName, movieTitle, showtime string, seats []string, amount float64) error
	}
}

func (a *emailAdapter) SendPaymentReceipt(toEmail, toName, movieTitle, showtime string, seats []string, amount float64) error {
	return a.sender.SendPaymentReceipt(toEmail, toName, movieTitle, showtime, seats, amount)
}
