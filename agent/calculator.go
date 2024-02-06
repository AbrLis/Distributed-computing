package agent

import (
	"sync"
	"time"
)

var done chan struct{} // Канал завершения вычислительных операций

// FreeCalculators - Структура счётчика свободных вычислителей
type FreeCalculators struct {
	Count           int         // Свободные вычислители
	PingTimeoutCalc []time.Time // Таймауты пингов вычислителей
	mu              sync.Mutex
}

// NewFreeCalculators создает новый экземпляр структуры счётчика свободных вычислителей
func NewFreeCalculators() *FreeCalculators {
	return &FreeCalculators{
		Count:           CountCalculators,
		PingTimeoutCalc: make([]time.Time, CountCalculators),
	}
}

// RunCalculators запускает вычислители ожидающие очередь задач
func (c *FreeCalculators) RunCalculators() {
	for i := 0; i < c.Count; i++ {
		go func(idCalc int) {
			for {
				select {
				case token := <-taskChannel:
					// TODO: выполнение вычислений
					_ = token
				case <-time.After(3 * time.Second): // Пингуемся записывая текущее время в PingTimeoutCalc
					c.mu.Lock()
					c.PingTimeoutCalc[idCalc] = time.Now()
					c.mu.Unlock()
				case <-done:
					// Завершение операций
					return
				}
			}
		}(i)
	}
}

// TODO: реализовать проверку на зависание задачи на основе пингов вычислителей
