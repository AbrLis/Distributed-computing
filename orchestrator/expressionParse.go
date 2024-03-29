package orchestrator

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/AbrLis/Distributed-computing/agent"
)

// ParseExpression разбивает выражение на токены в польской нотации
func ParseExpression(expr string) ([]agent.Token, error) {
	var tokens []agent.Token
	var ops []rune
	var buffer string

	for _, ch := range expr {
		if unicode.IsDigit(ch) {
			buffer += string(ch)
		} else if strings.ContainsRune("+-*/", ch) {
			if buffer != "" {
				tokens = append(tokens, agent.Token{Value: buffer, IsOp: false})
				buffer = ""
			}
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(ch) {
				tokens = append(tokens, agent.Token{Value: string(ops[len(ops)-1]), IsOp: true})
				ops = ops[:len(ops)-1]
			}
			ops = append(ops, ch)
		} else if ch != ' ' {
			return nil, fmt.Errorf("invalid character: %c", ch)
		}
	}

	if buffer != "" {
		tokens = append(tokens, agent.Token{Value: buffer, IsOp: false})
	}

	for len(ops) > 0 {
		tokens = append(tokens, agent.Token{Value: string(ops[len(ops)-1]), IsOp: true})
		ops = ops[:len(ops)-1]
	}

	return tokens, nil
}

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return 0
	}
}
