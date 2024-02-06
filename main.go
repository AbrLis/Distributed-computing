package main

import (
	"github.com/AbrLis/Distributed-computing/agent"
	"github.com/AbrLis/Distributed-computing/database"
	apiEndpoint "github.com/AbrLis/Distributed-computing/orchestrator"
)

func main() {
	calculators := agent.NewFreeCalculators()
	calculators.RunCalculators() // Запуск вычислительных операций
	// TODO: Запуск демона RunDeamon(calculators)

	db := database.NewDatabase()

	orchestrator := apiEndpoint.NewOrchestrator(db)
	err := orchestrator.Run("http://localhost", "3000")
	if err != nil {
		panic(err)
	}
}
