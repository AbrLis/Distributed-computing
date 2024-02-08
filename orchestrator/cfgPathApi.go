package orchestrator

import (
	"github.com/AbrLis/Distributed-computing/database"
	"time"
)

// Константы для путей API
const (
	addExpressionPath  = "/add-expression"  // API для добавления арифметического выражения
	getExpressionsPath = "/get-expressions" // API для получения списка арифметических выражений
	getValuePath       = "/get-value/"      // API для получения значения выражения по его идентификатору
	getOperationsPath  = "/get-operations"  // API для получения списка доступных операций со временем их выполнения
	GetTaskPath        = "/get-task"        // API для получения задачи для выполнения
	ReceiveResultPath  = "/receive-result"  // API для приема результата обработки данных
	HostPath           = "localhost"        // Путь до хоста
	PortHost           = ":3000"            // Порт хоста
)

// OpetatorTimeout - структура для формирования списка доступных операций со временем их выполнения
type OpetatorTimeout struct {
	Add  string `json:"+"`
	Sub  string `json:"-"`
	Mult string `json:"*"`
	Div  string `json:"/"`
}

// ExpressionStatus - структура для формирования статуса выражения
type ExpressionStatus struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Status     string `json:"status"`
	Result     string `json:"result,omitempty"`
}

// SendREsult - структура отправки результата оркестратору
type SendREsult struct {
	IDCalc int64               `json:"id"`
	Result string              `json:"result"`
	Status database.TaskStatus `json:"status"`
}

// Token - структура для формирования польской нотации выражения
type Token struct {
	Value string
	IsOp  bool
}

// TaskCalculate - структура для формирования задачи
type TaskCalculate struct {
	ID         int64   `json:"id"`
	Expression []Token `json:"expression"`
}

// ========================================
// Константы таймаутов для операций +, -, *
const (
	AddTimeout      = 5 * time.Second
	SubtractTimeout = 3 * time.Second
	MultiplyTimeout = 4 * time.Second
	DivideTimeout   = 6 * time.Second
)

const (
	CountCalculators      = 5                // Количество вычислителей в демоне
	InactiveServerTimeout = 10 * time.Second // Таймаут после которого вычислитель будет считаться неактивным
)

// ========================================
