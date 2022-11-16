package mail

import (
	"encoding/base64"
	"net/mail"
	"net/smtp"

	"github.com/alexm24/golang/internal/models"
)

type Mail struct {
}

func NewMail() *Mail {
	return &Mail{}
}

func (m *Mail) SendMail(item models.Zoom) error {
	c, err := smtp.Dial("10.0.16.1:25")
	if err != nil {
		return err
	}
	fromEmail := "null@vp.ru"
	from := (&mail.Address{Name: "Запись Zoom", Address: fromEmail}).String()

	if err = c.Mail(fromEmail); err != nil {
		return err
	}

	if err = c.Rcpt(*item.Email); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	body := "<h2>Запись zoom конференции находится по адресу: </h2>"
	body += "<h2><a href=\"https://vp.ru/zoom/" + item.Id.String() + "\">" + *item.Topic + "</a></h2>"

	msg := "From: " + from + "\r\n" +
		"Subject: Запись Zoom" + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}

	return nil
}
