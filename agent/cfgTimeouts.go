package agent

import (
	"time"
)

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
