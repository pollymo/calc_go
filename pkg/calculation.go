package pkg

import (
	"fmt"
	"strconv"
)

type TokenType int

const (
	TOKEN_NUMBER TokenType = iota
	TOKEN_OPERATOR
)

type Token struct {
	Type  TokenType
	Value string
}

func isOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/'
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func infixToPostfix(infix string) []Token {
	var output []Token
	var operators []string

	for _, char := range infix {
		token := string(char)
		if isOperator(rune(char)) {
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				output = append(output, Token{TOKEN_OPERATOR, operators[len(operators)-1]})
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for operators[len(operators)-1] != "(" {
				output = append(output, Token{TOKEN_OPERATOR, operators[len(operators)-1]})
				operators = operators[:len(operators)-1]
			}
			operators = operators[:len(operators)-1]
		} else {
			output = append(output, Token{TOKEN_NUMBER, token})
		}
	}

	for len(operators) > 0 {
		output = append(output, Token{TOKEN_OPERATOR, operators[len(operators)-1]})
		operators = operators[:len(operators)-1]
	}

	return output
}

func evaluatePostfix(postfix []Token) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if token.Type == TOKEN_NUMBER {
			num, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("invalid expression")
			}
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			var result float64
			switch token.Value {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				result = a / b
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}
	return stack[0], nil
}

func Calc(expression string) (float64, error) {
	postfix := infixToPostfix(expression)
	return evaluatePostfix(postfix)
}
