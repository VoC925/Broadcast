package client

import "fmt"

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
	fmt.Printf("data sent to number ... %s\n", s.number)
	// maybe realisation ...
	return nil
}

// TODO: Сделать реализацию отправки на почту
