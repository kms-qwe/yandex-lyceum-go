package orchestrator

import (
	"container/list"
	"errors"
	"strings"
)

// isOperator проверяет, является ли символ оператором
func isOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '(' || c == ')'
}

// precedence возвращает приоритет оператора
func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	case '(', ')':
		return 3
	}
	return 0
}

// infixToPostfix переводит выражение из инфиксной нотации в обратную польскую
func infixToPostfix(infix string) (string, error) {
	var postfix strings.Builder
	stack := list.New()

	for _, char := range infix {
		if char == ' ' {
			continue
		} else if char >= '0' && char <= '9' {
			postfix.WriteRune(char)
		} else if isOperator(char) {
			if char == '(' {
				stack.PushBack(char)
			} else if char == ')' {
				for stack.Len() > 0 {
					top := stack.Remove(stack.Back()).(rune)
					if top != '(' {
						postfix.WriteRune(top)
					} else {
						break
					}
				}
			} else {
				for stack.Len() > 0 {
					top := stack.Back().Value.(rune)
					if precedence(top) >= precedence(char) {
						postfix.WriteRune(top)
						stack.Remove(stack.Back())
					} else {
						break
					}
				}
				stack.PushBack(char)
			}
		} else {
			return "", errors.New("invalid data")
		}
	}

	for stack.Len() > 0 {
		postfix.WriteRune(stack.Remove(stack.Back()).(rune))
	}

	return postfix.String(), nil
}
