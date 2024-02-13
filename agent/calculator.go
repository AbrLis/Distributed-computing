package agent

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/AbrLis/Distributed-computing/database"
)

var done chan struct{} // Канал завершения вычислительных операций

// Token - структура для формирования польской нотации выражения
type Token struct {
	Value string
	IsOp  bool
}

// TaskCalculate - структура для формирования задачи
type TaskCalculate struct {
	ID         string  `json:"id"`
	Expression []Token `json:"expression"`
}

// FreeCalculators - Структура счётчика свободных вычислителей
type FreeCalculators struct {
	db              *database.Database       // Ссылка на бд
	Count           int                      // Количество вычислителей
	CountFree       int                      // Свободные вычислители
	PingTimeoutCalc []time.Time              // Таймауты пингов вычислителей
	Queue           []TaskCalculate          // Очередь исполнения задач
	queueInProcess  map[string]TaskCalculate // Задачи находящиеся на обработке
	taskChannel     chan TaskCalculate       // Канал задач
	AddTimeout      time.Duration            // Таймауты операций
	SubtractTimeout time.Duration
	MultiplyTimeout time.Duration
	DivideTimeout   time.Duration
	mu              sync.Mutex
}

// NewFreeCalculators создает новый экземпляр структуры счётчика свободных вычислителей
func NewFreeCalculators(db *database.Database) *FreeCalculators {
	freeCaclulatros := &FreeCalculators{
		db:              db,
		Count:           5,
		CountFree:       5,
		Queue:           []TaskCalculate{},
		queueInProcess:  map[string]TaskCalculate{},
		taskChannel:     make(chan TaskCalculate),
		AddTimeout:      5 * time.Second,
		SubtractTimeout: 3 * time.Second,
		MultiplyTimeout: 4 * time.Second,
		DivideTimeout:   6 * time.Second,
		mu:              sync.Mutex{},
	}
	freeCaclulatros.PingTimeoutCalc = make([]time.Time, freeCaclulatros.Count)

	return freeCaclulatros
}

// RunCalculators запускает вычислители ожидающие очередь задач
func (c *FreeCalculators) RunCalculators() {
	for i := 0; i < c.Count; i++ {
		go func(calcId int) {
			for {
				select {
				case tokens := <-c.taskChannel:
					log.Printf("Вычислитель %d - получил задачу: %s\n", calcId, tokens.ID)
					result, flagError := c.calculateValue(calcId, tokens.Expression)

					c.sendResult(tokens.ID, flagError, result)
					log.Println("Вычислитель отправил результат в бд: ", tokens.ID)

					// Переход в режим ожидания
					c.CountFree++
					continue

				case <-time.After(3 * time.Second): // Пингуемся записывая текущее время в PingTimeoutCalc
					log.Println("Вычислитель пингуется: ", calcId)
					c.mu.Lock()
					c.PingTimeoutCalc[calcId] = time.Now()
					c.mu.Unlock()
				case <-done:
					// Завершение операций
					log.Println("Вычислитель завершил работу: ", calcId)
					return
				}
			}
		}(i)
	}
}

// sendResult - Отправка результата на оркестратор
func (c *FreeCalculators) sendResult(idCalc string, flagError bool, result float64) {
	// Отправка результатов и переход в режим ожидания
	textResult := "error parse or calculate"
	status := database.StatusError
	if !flagError {
		textResult = strconv.FormatFloat(result, 'f', -1, 64)
		status = database.StatusCompleted
	}

	// Запись результатов обработки в базу данных
	err := c.db.SetTaskResult(idCalc, status, textResult)
	if err != nil {
		log.Printf("Ошибка записи результата в базу данных: %s\n", err)
		c.mu.Lock()
		c.Queue = append(c.Queue, c.queueInProcess[idCalc])
		c.mu.Unlock()
		log.Println("Ошибочная операция перенесена в конец очереди...")
	}

	// Удаление задачи из очереди обработки
	delete(c.queueInProcess, idCalc)
}

// calculateValue вычисляет значение выражения
func (c *FreeCalculators) calculateValue(idCalc int, tokens []Token) (float64, bool) {
	var result float64
	flagError := false // Признак ошибки при выполнении операции
	if len(tokens) == 0 {
		flagError = true
		log.Println("Очередь задач пустая, вычисление невозможно")
		flagError = true
	} else {
		// Вычисление выражения
		stack := make([]float64, 0)
		for _, token := range tokens {
			if !token.IsOp {
				num, err := strconv.ParseFloat(token.Value, 64)
				if err != nil {
					log.Println("Ошибка при парсинге числа в вычислителе", err)
					flagError = true
					break
				}
				stack = append(stack, num)
			} else {
				if len(stack) < 2 {
					log.Println("Для операции необходимо два числа в стеке, ошибка в вычислителе")
					flagError = true
					break
				}
				num1, num2 := stack[len(stack)-2], stack[len(stack)-1]
				stack = stack[:len(stack)-2]

				switch token.Value {
				case "+":
					stack = append(stack, num1+num2)
					time.Sleep(c.AddTimeout)
				case "-":
					stack = append(stack, num1-num2)
					time.Sleep(c.SubtractTimeout)
				case "*":
					stack = append(stack, num1*num2)
					time.Sleep(c.MultiplyTimeout)
				case "/":
					if num2 == 0 {
						log.Println("Деление на ноль")
						flagError = true
						break
					}
					stack = append(stack, num1/num2)
					time.Sleep(c.DivideTimeout)
				default:
					log.Println("Неизвестная операция в вычислителе")
					flagError = true
					break
				}

				// После кажого вычисления отправка пинга что вычислитель жив
				c.mu.Lock()
				c.PingTimeoutCalc[idCalc] = time.Now()
				c.mu.Unlock()
			}
		}
		if len(stack) != 1 {
			log.Println("Слишком много чисел в стеке, ошибка в вычислителе")
			flagError = true
		}
		result = stack[0]
	}
	return result, flagError
}

// TODO: реализовать проверку на зависание задачи на основе пингов вычислителей
