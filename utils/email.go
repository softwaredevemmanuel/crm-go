package utils

import (
	"gopkg.in/gomail.v2"
	"log"
	"crm-go/config"
)

var cfg = config.LoadConfig()

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(cfg.SMTPFrom, "Ehizua Hub Learning Center"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(cfg.SMTPServer, cfg.SMTPPort, cfg.SMTPLogin, cfg.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("❌ Failed to send email: %v", err)
		return err
	}
	log.Println("📨 Email sent successfully to", to)
	return nil
}
