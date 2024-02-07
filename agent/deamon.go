package agent

import (
	"github.com/AbrLis/Distributed-computing/orchestrator"
)

var taskChannel chan []orchestrator.Token // Канал задач

// TODO: если есть свободные вычислители, то запрос у оркестратора очереди вычислений, передать в канал эту очередь на исполнение.
