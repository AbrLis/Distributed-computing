package agent

import (
	"log"
	"time"
)

// RunDeamon запускает демона вычислителя ожидающего задачи
func RunDeamon(calculator *FreeCalculators) {
	for {
		if calculator.CountFree > 0 {
			calculator.mu.Lock()
			if len(calculator.Queue) == 0 {
				calculator.mu.Unlock()
				log.Println("Задач в очереди нет, ожидание 2 секунды...")
				time.Sleep(2 * time.Second)
				continue
			}

			// Отправка задачи свободным вычислителям
			calculator.CountFree--
			task := calculator.Queue[0]
			// Перевод из очереди ожидания в очередь обработки
			calculator.queueInProcess[calculator.Queue[0].ID] = calculator.Queue[0]
			calculator.Queue = calculator.Queue[1:]
			calculator.mu.Unlock()
			calculator.taskChannel <- task
		} else {
			// Ждать пока не появятся свободные вычислители
			log.Println("Нет свободных вычислителей, ожидание...")
			time.Sleep(2 * time.Second)
			continue
		}
	}
}
