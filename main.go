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
	log.Println("Попытка запуска API-оркестратора")
	err := orchestrator.Run() // Запуск определён по адресу localhost:3000
	if err != nil {
		log.Println("Ошибка запуска API-оркестратора")
		panic(err)
	}
	log.Printf("API-оркестратор запущен по адресу http://%s:%s\n", apiEndpoint.HostPath, apiEndpoint.PortHost)

	calculators := agent.NewFreeCalculators()
	calculators.RunCalculators() // Запуск вычислительных операций
	log.Println("Вычислители запущены")
	agent.RunDeamon(calculators) // Запуск демона вычислителя
	log.Println("Демон вычислитетей запущен")
	log.Println("Система готова к использованию")
}
