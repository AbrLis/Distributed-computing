package main

import (
	"github.com/AbrLis/Distributed-computing/agent"
	"github.com/AbrLis/Distributed-computing/database"
	apiEndpoint "github.com/AbrLis/Distributed-computing/orchestrator"
)

func main() {
	db := database.NewDatabase()

	orchestrator := apiEndpoint.NewOrchestrator(db)
	err := orchestrator.Run() // Запуск определён по адресу localhost:3000
	if err != nil {
		panic(err)
	}

	calculators := agent.NewFreeCalculators()
	calculators.RunCalculators() // Запуск вычислительных операций
	agent.RunDeamon(calculators)
}
