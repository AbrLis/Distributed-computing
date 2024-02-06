package orchestrator

// Константы для путей API
const (
	addExpressionPath  = "/add-expression"  // API для добавления арифметического выражения
	getExpressionsPath = "/get-expressions" // API для получения списка арифметических выражений
	getValuePath       = "/get-value/"      // API для получения значения выражения по его идентификатору
	getOperationsPath  = "/get-operations"  // API для получения списка доступных операций со временем их выполнения
	getTaskPath        = "/get-task"        // API для получения задачи для выполнения
	receiveResultPath  = "/receive-result"  // API для приема результата обработки данных
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
