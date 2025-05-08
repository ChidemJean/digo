package main

import (
	"fmt"
	"log"

	"github.com/ChidemJean/digo"
)

type Notificador interface {
	Notificar(msg string)
}

type EmailNotificador struct{}

func (e *EmailNotificador) Notificar(msg string) {
	fmt.Println("Email enviado:", msg)
}

func NewEmailNotificador() *EmailNotificador {
	return &EmailNotificador{}
}

type AlertaService struct {
	Notificador Notificador
}

func NewAlertaService(n Notificador) *AlertaService {
	return &AlertaService{Notificador: n}
}

func (a *AlertaService) Enviar() {
	a.Notificador.Notificar("Sistema em alerta!")
}

func main() {
	container := digo.New()
	container.Register(NewEmailNotificador, digo.Singleton)
	container.Register(NewAlertaService, digo.Transient)
	container.RegisterInterface((*Notificador)(nil), &EmailNotificador{})

	alerta, err := digo.Resolve[Notificador]()
	if err != nil {
		log.Fatal(err)
	}
	alerta.Notificar("Teste")
}
