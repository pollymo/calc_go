package calculation

import (
	"strconv"
	"unicode"
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

func infixToPostfix(infix string) ([]Token, error) {
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
		} else if unicode.IsDigit(rune(char)) {
			output = append(output, Token{TOKEN_NUMBER, token})
		} else {
			return nil, ErrUnexpectedToken
		}
	}

	for len(operators) > 0 {
		output = append(output, Token{TOKEN_OPERATOR, operators[len(operators)-1]})
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

func evaluatePostfix(postfix []Token) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if token.Type == TOKEN_NUMBER {
			num, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, ErrInvalidNumber
			}
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, ErrNotEnoughValues
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
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				result = a / b
			default:
				return 0, ErrInvalidOperator
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}
	return stack[0], nil
}

func Calc(expression string) (float64, error) {
	if expression == "" {
		return 0, ErrEmptyInput
	}
	postfix, _ := infixToPostfix(expression)
	return evaluatePostfix(postfix)
}
