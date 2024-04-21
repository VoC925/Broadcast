package client

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// интерфейс получателя прогноза погоды
type Sender interface {
	Send(data []byte) error
}

// Структура, реализующая интерфейс
type SenderToEmail struct {
	number string
}

// конструктор
func NewSenderToEmail(number string) *SenderToEmail {
	return &SenderToEmail{
		number: number,
	}
}

func (s *SenderToEmail) Send(data []byte) error {
	// fmt.Printf("data sent to number ... %s\n", s.number)
	logrus.WithFields(logrus.Fields{
		"msg": string(data),
	}).Info("data sent")
	// maybe realisation ...
	return nil
}

type HermesSender struct {
	hermes *hermes.Hermes  // hermes для создания красивых сообщений HTML
	sender *gomail.Message // gomail для непосредственной отправки сообщения из hermes
}

// конструктор HermesSender, где входные параметры:
// email - почта, на которую должны приходить прогнозы погоды
// themeName - тема сообщения.
func NewHermesSender(email, themeName string) *HermesSender {
	return &HermesSender{
		hermes: &hermes.Hermes{
			Theme: new(hermes.Default),
		},
		sender: gomail.NewMessage(),
	}
}

func (hs *HermesSender) Send(data []byte) error {
	// экземпляр сообщения
	email := hermes.Email{
		Body: hermes.Body{
			Name: "Eugene", // приветствие в начале сообщения: "Hi, Eugene"
			Intros: []string{
				fmt.Sprintf("Weather in your location for 1 hour: %s °C", string(data)),
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	// генерация HTML на основе сообщения
	emailBody, err := hs.hermes.GenerateHTML(email)
	if err != nil {
		return err
	}

	// создание электронного письма
	hs.sender.SetHeader("From", "ewg.covaleov1999@yandex.ru")
	hs.sender.SetHeader("To", "ewg.covaleov1999@yandex.ru")
	hs.sender.SetHeader("Subject", "Broadcast")
	hs.sender.AddAlternative("text/html", emailBody)

	d := gomail.NewDialer("smtp.yandex.ru", 587, "ewg.covaleov1999@yandex.ru", "kov999ALEV2039")

	if err := d.DialAndSend(hs.sender); err != nil {
		return err
	}

	return nil
}
