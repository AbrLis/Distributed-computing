package orchestrator

// OpetatorTimeout - структура для формирования списка доступных операций со временем их выполнения
type OpetatorTimeout struct {
	Add  string `json:"add"`
	Sub  string `json:"sub"`
	Mult string `json:"mult"`
	Div  string `json:"div"`
}

// ExpressionStatus - структура для формирования статуса выражения
type ExpressionStatus struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Status     string `json:"status"`
	Result     string `json:"result,omitempty"`
}
