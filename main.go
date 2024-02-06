package main

import (
	"fmt"
	"github.com/AbrLis/Distributed-computing/orchestrator"
)

func main() {
	expr := "22+22*4"
	tokens, err := orchestrator.ParseExpression(expr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, token := range tokens {
		if token.IsOp {
			fmt.Println("Operator:", token.Value)
		} else {
			fmt.Println("Operand:", token.Value)
		}
	}
}
