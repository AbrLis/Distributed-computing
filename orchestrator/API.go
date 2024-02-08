// В файле orchestrator.go

package orchestrator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AbrLis/Distributed-computing/database"
)

// Orchestrator представляет сервер-оркестратор
type Orchestrator struct {
	db             *database.Database
	queue          [][]Token         // Очередь исполнения задач
	queueInProcess map[int64][]Token // Задачи находящиеся на обработке
}

// NewOrchestrator создает новый экземпляр оркестратора
func NewOrchestrator(db *database.Database) *Orchestrator {
	return &Orchestrator{
		db:             db,
		queue:          [][]Token{},
		queueInProcess: map[int64][]Token{},
	}
}

// AddExpressionHandler обработчик для добавления арифметического выражения
func (o *Orchestrator) AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Не удалось прочитать тело запроса", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Проверяем, что тело запроса не пустое
	if len(body) == 0 {
		http.Error(w, "Тело запроса пустое", http.StatusBadRequest)
		return
	}

	// Преобразуем тело запроса в строку
	expression := string(body)

	// Парсим выражение
	tokens, err := ParseExpression(expression)
	if err != nil {
		http.Error(w, "Неверный формат выражения", http.StatusBadRequest)
		return
	}

	// Добавляем выражение в базу данных
	id := o.db.GenerateID()
	o.db.AddTask(id, expression)

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Выражение добавлено в базу данных и принято к обработке. ID: " + id))

	// Добавление задачи в очередь
	o.queue = append(o.queue, tokens)

}

// GetExpressionsHandler обработчик для получения списка выражений со статусами
func (o *Orchestrator) GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем все задачи из базы данных
	tasks := o.db.GetAllTasks()

	// Формируем список выражений со статусами
	expressions := make([]ExpressionStatus, 0)
	for id, task := range tasks {
		status := database.GetTaskStatus(task.Status)

		expressions = append(
			expressions, ExpressionStatus{
				ID:         id,
				Expression: task.Expression,
				Status:     status,
				Result:     task.Result,
			},
		)
	}

	// Кодируем список выражений в JSON и отправляем клиенту
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(expressions); err != nil {
		http.Error(w, "Ошибка при кодировании данных в JSON - выражений", http.StatusInternalServerError)
		return
	}
}

// GetValueHandler обработчик для получения значения выражения по его идентификатору
func (o *Orchestrator) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем идентификатор задачи из URL
	id := strings.TrimPrefix(r.URL.Path, getValuePath)

	// Получаем задачу по идентификатору из базы данных
	task, exists := o.db.GetTask(id)
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Проверяем статус задачи
	switch task.Status {
	case database.StatusInProgress:
		http.Error(w, "Задача находится в процессе выполнения", http.StatusConflict)
		return
	case database.StatusCompleted:
		// Отправляем результат вычисления выражения
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Результат вычисления выражения: " + task.Result))
	default:
		http.Error(w, "Задача находится в невалидном статусе", http.StatusInternalServerError)
	}
}

// GetOperationsHandler обработчик для получения списка доступных операций со временем их выполнения
func (o *Orchestrator) GetOperationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	result := OpetatorTimeout{
		Add:  AddTimeout.String(),
		Sub:  SubtractTimeout.String(),
		Mult: MultiplyTimeout.String(),
		Div:  DivideTimeout.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Ошибка при кодировании данных в JSON - таймаутов", http.StatusInternalServerError)
	}
}

// GetTaskHandler обработчик для получения задачи для выполнения
func (o *Orchestrator) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	if o.queue == nil || len(o.queue) == 0 {
		http.Error(w, "Очередь задач пустая", http.StatusNotFound)
		return
	}

	id := time.Now().UnixNano()
	taskCalculate := TaskCalculate{
		ID:         id,
		Expression: o.queue[0],
	}

	// Отправка задачи демону распределителю
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(taskCalculate); err != nil {
		errorsText := "Ошибка при кодировании данных в JSON - задачи"
		http.Error(w, errorsText, http.StatusInternalServerError)
		log.Println(errorsText)
		return
	}

	// Перевод из очереди ожидания в очередь обработки
	o.queueInProcess[id] = o.queue[0]
	o.queue = o.queue[1:]
}

// ReceiveResultHandler обработчик для приема результата обработки данных
func (o *Orchestrator) ReceiveResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	result := &SendREsult{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		http.Error(w, "Не удалось распарсить тело запроса", http.StatusBadRequest)
		return
	}

	// Запись результатов обработки в базу данных
	err = o.db.SetTaskResult(strconv.FormatInt(result.IDCalc, 10), result.Status, result.Result)
	if err != nil {
		http.Error(w, "Ошибка при записи результатов в базу данных", http.StatusInternalServerError)
		return
	}

	// Удаление задачи из очереди обработки
	delete(o.queueInProcess, int64(result.IDCalc))

	w.WriteHeader(http.StatusOK)
}

// Run запускает сервер оркестратора
func (o *Orchestrator) Run() error {
	http.HandleFunc(addExpressionPath, o.AddExpressionHandler)
	http.HandleFunc(getExpressionsPath, o.GetExpressionsHandler)
	http.HandleFunc(getValuePath, o.GetValueHandler)
	http.HandleFunc(getOperationsPath, o.GetOperationsHandler)
	http.HandleFunc(GetTaskPath, o.GetTaskHandler)
	http.HandleFunc(ReceiveResultPath, o.ReceiveResultHandler)

	errCh := make(chan error)
	go func() {
		errCh <- http.ListenAndServe(HostPath+PortHost, nil)
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(2 * time.Second):
		return nil
	}
}
