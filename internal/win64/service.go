package win64

import (
	"fmt"
	"src/post_relay/internal/db"

	"github.com/kardianos/service"
)

type Program struct{}

// Start define o que o serviço vai fazer quando iniciado
func (p *Program) Start(s service.Service) error {
	// A lógica do seu serviço começa aqui
	go p.run()
	return nil
}

// run contém a lógica que o serviço vai executar enquanto estiver rodando
func (p *Program) run() {
	db.StartNotifications()
}

// Stop define o que acontece quando o serviço é parado
func (p *Program) Stop(s service.Service) error {
	fmt.Println("O serviço foi parado.")
	return nil
}

// NewService cria uma nova instância do serviço
func NewService() (service.Service, error) {
	prg := &Program{}
	svcConfig := &service.Config{
		Name:        "Attom",                        // Nome do serviço
		DisplayName: "Attom",                        // Nome para exibição no Gerenciador de Serviços
		Description: "Este é um serviço Go simples", // Descrição do serviço
	}

	// Cria e retorna a instância do serviço
	return service.New(prg, svcConfig)
}
