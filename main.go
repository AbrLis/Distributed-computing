package main

import (
	"log"

	"github.com/AbrLis/Distributed-computing/agent"
	"github.com/AbrLis/Distributed-computing/database"
	apiEndpoint "github.com/AbrLis/Distributed-computing/orchestrator"
)

func main() {
	db := database.NewDatabase()

	orchestrator := apiEndpoint.NewOrchestrator(db)

	// Запуск API-оркестратора
	err := orchestrator.Run() // Запуск определён по адресу localhost:3000
	if err != nil {
		log.Println("Ошибка запуска API-оркестратора")
		panic(err)
	}

	calculators := agent.NewFreeCalculators()
	calculators.RunCalculators() // Запуск вычислительных операций
	agent.RunDeamon(calculators) // Запуск демона вычислителя
}
