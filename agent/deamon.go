package agent

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/AbrLis/Distributed-computing/orchestrator"
)

// RunDeamon запускает демона вычислителя ожидающего задачи
func RunDeamon(calculator *FreeCalculators) {
	for {
		// Попытка запросить задачу
		if calculator.Count > 0 {
			// TODO: запросить задачу у оркестратора
			req, err := http.NewRequest(
				http.MethodGet, "http://"+orchestrator.HostPath+orchestrator.PortHost+orchestrator.GetTaskPath, nil,
			)
			if err != nil {
				log.Println("Не удалось создать запрос в оркестратор")
				time.Sleep(1 * time.Second)
				continue
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("Ошибка при посылке запроса в оркестратор")
				time.Sleep(1 * time.Second)
				continue
			}
			if resp.StatusCode != 200 {
				log.Println("Задач в очереди нет, ожидание 2 секунды...")
				time.Sleep(2 * time.Second)
				continue
			}
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Println("Не удалось прочитать тело ответа")
				time.Sleep(1 * time.Second)
				continue
			}

			// Попытка преобразовать ответ в очередь вычислительных задач
			tokens := orchestrator.TaskCalculate{}
			err = json.Unmarshal(body, &tokens)
			if err != nil {
				log.Println("Не удалось преобразовать тело ответа, json Unmarshal")
				time.Sleep(1 * time.Second)
				continue
			}

			// Задача получена, отправка свободным вычислителям
			calculator.Count--
			calculator.taskChannel <- tokens

		} else {
			// Ждать пока не появятся свободные вычислители
			log.Println("Нет свободных вычислителей, ожидание...")
			time.Sleep(2 * time.Second)
			continue
		}
	}
}
