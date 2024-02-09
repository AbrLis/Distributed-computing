package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/AbrLis/Distributed-computing/database"
	"github.com/AbrLis/Distributed-computing/orchestrator"
)

var done chan struct{} // Канал завершения вычислительных операций
//var taskChannel chan orchestrator.TaskCalculate // Канал задач

// FreeCalculators - Структура счётчика свободных вычислителей
type FreeCalculators struct {
	Count           int                             // Свободные вычислители
	PingTimeoutCalc []time.Time                     // Таймауты пингов вычислителей
	taskChannel     chan orchestrator.TaskCalculate // Канал задач
	mu              sync.Mutex
}

// NewFreeCalculators создает новый экземпляр структуры счётчика свободных вычислителей
func NewFreeCalculators() *FreeCalculators {
	return &FreeCalculators{
		Count:           orchestrator.CountCalculators,
		PingTimeoutCalc: make([]time.Time, orchestrator.CountCalculators),
		taskChannel:     make(chan orchestrator.TaskCalculate),
	}
}

// RunCalculators запускает вычислители ожидающие очередь задач
func (c *FreeCalculators) RunCalculators() {
	for i := 0; i < c.Count; i++ {
		go func(calcId int) {
			for {
				select {
				case tokens := <-c.taskChannel:
					log.Println("Вычислитель получил задачу: ", tokens.ID)
					result, flagError := c.calculateValue(calcId, tokens.Expression)

					c.sendResult(tokens.ID, flagError, result)

					// Переход в режим ожидания
					c.Count++
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

	sendresult := orchestrator.SendREsult{
		IDCalc: idCalc,
		Result: textResult,
		Status: status,
	}

	// TODO: Безобразное игнорирование всех ошибок для простоты, иначе боюсь не успею.
	jsonResult, _ := json.Marshal(sendresult)
	req, _ := http.NewRequest(
		"POST",
		"http://"+orchestrator.HostPath+orchestrator.PortHost+orchestrator.ReceiveResultPath,
		bytes.NewBuffer(jsonResult),
	)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		log.Printf("Не удалось отправить результат вычисления выражения в оркестратор id:%d\n", idCalc)
	}
}

// calculateValue вычисляет значение выражения
func (c *FreeCalculators) calculateValue(idCalc int, tokens []orchestrator.Token) (float64, bool) {
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
					time.Sleep(orchestrator.AddTimeout)
				case "-":
					stack = append(stack, num1-num2)
					time.Sleep(orchestrator.SubtractTimeout)
				case "*":
					stack = append(stack, num1*num2)
					time.Sleep(orchestrator.MultiplyTimeout)
				case "/":
					if num2 == 0 {
						log.Println("Деление на ноль")
						flagError = true
						break
					}
					stack = append(stack, num1/num2)
					time.Sleep(orchestrator.DivideTimeout)
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
