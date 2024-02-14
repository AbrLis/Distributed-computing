package orchestrator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AbrLis/Distributed-computing/agent"
	"github.com/AbrLis/Distributed-computing/database"
)

// Orchestrator представляет сервер-оркестратор
type Orchestrator struct {
	db         *database.Database
	calculator *agent.FreeCalculators
}

// NewOrchestrator создает новый экземпляр оркестратора
func NewOrchestrator(db *database.Database, calc *agent.FreeCalculators) *Orchestrator {
	return &Orchestrator{
		db:         db,
		calculator: calc,
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
	o.calculator.Queue = append(o.calculator.Queue, agent.TaskCalculate{ID: id, Expression: tokens})

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

	// Формирование результата вычисления выражения
	expression := ExpressionStatus{
		ID:         id,
		Expression: task.Expression,
		Status:     database.GetTaskStatus(task.Status),
		Result:     task.Result,
	}

	// Кодируем список выражений в JSON и отправляем клиенту
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(expression); err != nil {
		http.Error(w, "Ошибка при кодировании данных в JSON - выражения", http.StatusInternalServerError)
		return
	}
}

// GetOperationsHandler обработчик для получения списка доступных операций со временем их выполнения
func (o *Orchestrator) GetOperationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	result := OpetatorTimeout{
		Add:  o.calculator.AddTimeout.String(),
		Sub:  o.calculator.SubtractTimeout.String(),
		Mult: o.calculator.MultiplyTimeout.String(),
		Div:  o.calculator.DivideTimeout.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Ошибка при кодировании данных в JSON - таймаутов", http.StatusInternalServerError)
	}
}

// MonitoringHandler обработчик для получения статуса вычислителей
func (o *Orchestrator) MonitoringHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	pingTimeout := make(map[int]string)
	for i, value := range o.calculator.PingTimeoutCalc {
		pingTimeout[i+1] = fmt.Sprintf("%.3f", time.Now().Sub(value).Seconds()) + " sec"
	}

	jsonData, err := json.Marshal(pingTimeout)
	if err != nil {
		http.Error(w, "Ошибка при кодировании данных в JSON - статуса вычислителя", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// Run запускает сервер оркестратора
func (o *Orchestrator) Run() error {
	http.HandleFunc(addExpressionPath, o.AddExpressionHandler)
	http.HandleFunc(getExpressionsPath, o.GetExpressionsHandler)
	http.HandleFunc(getValuePath, o.GetValueHandler)
	http.HandleFunc(getOperationsPath, o.GetOperationsHandler)
	http.HandleFunc(monitoring, o.MonitoringHandler)

	errCh := make(chan error)
	go func() {
		errCh <- http.ListenAndServe(HostPath+PortHost, nil)
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(2 * time.Second):
		close(errCh)
		return nil
	}
}
