package orchestrator

import (
	"container/list"
	"errors"
	"strings"
	"unicode"
)

func isOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '(' || c == ')'
}

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	case '(':
		return 0
	}
	return -1
}

func isValidExpression(expression string) bool {
	if len(expression) == 0 {
		return false
	}

	// Стек для отслеживания скобок
	var stack []rune
	lastChar := ' ' // Переменная для хранения предыдущего символа

	for i, char := range expression {
		switch {
		case unicode.IsDigit(char):
			// Если символ - цифра, продолжаем
			lastChar = char

		case char == '+', char == '-', char == '*', char == '/':
			// Проверяем, является ли текущий символ унарным минусом или плюсом
			if char == '+' || char == '-' {
				if i == 0 || expression[i-1] == '(' || expression[i-1] == '+' || expression[i-1] == '-' || expression[i-1] == '*' || expression[i-1] == '/' {
					lastChar = char
					continue
				}
			}

			// Если символ - оператор, предыдущий символ не должен быть оператором или открывающей скобкой
			if lastChar == ' ' || lastChar == '+' || lastChar == '-' || lastChar == '*' || lastChar == '/' || lastChar == '(' {
				return false
			}
			lastChar = char

		case char == '(':
			// Если символ - открывающая скобка, добавляем её в стек
			stack = append(stack, char)
			lastChar = char

		case char == ')':
			// Если символ - закрывающая скобка, проверяем, есть ли соответствующая открывающая скобка в стеке
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return false
			}
			stack = stack[:len(stack)-1]
			lastChar = char

		case char == ' ':
			// Пропускаем пробелы
			continue

		default:
			// Если символ не является допустимым, возвращаем false
			return false
		}
	}

	// Проверяем, что стек пуст (все скобки закрыты) и последний символ не оператор
	return len(stack) == 0 && lastChar != '+' && lastChar != '-' && lastChar != '*' && lastChar != '/'
}

func infixToPostfix(infix string) (string, error) {
	if !isValidExpression(infix) {
		return "", errors.New("invalid string")
	}
	var postfix strings.Builder
	stack := list.New()
	var num strings.Builder
	isUnary := true

	for i, char := range infix {
		if unicode.IsSpace(char) {
			continue
		} else if unicode.IsDigit(char) {
			num.WriteRune(char)
			isUnary = false
		} else if isOperator(char) {
			if char == '/' {
				// Check for division by zero
				nextNonSpaceIdx := i + 1
				for nextNonSpaceIdx < len(infix) && unicode.IsSpace(rune(infix[nextNonSpaceIdx])) {
					nextNonSpaceIdx++
				}
				if nextNonSpaceIdx < len(infix) && rune(infix[nextNonSpaceIdx]) == '0' {
					return "", errors.New("division by zero")
				}
			}

			if num.Len() > 0 {
				postfix.WriteString(num.String() + " ")
				num.Reset()
			}

			if char == '(' {
				stack.PushBack(char)
				isUnary = true
			} else if char == ')' {
				for stack.Len() > 0 {
					top := stack.Remove(stack.Back()).(rune)
					if top == '(' {
						break
					}
					postfix.WriteRune(top)
					postfix.WriteRune(' ')
				}
				isUnary = false
			} else {
				if char == '-' && isUnary {
					num.WriteRune('-')
					isUnary = false
				} else {
					for stack.Len() > 0 {
						top := stack.Back().Value.(rune)
						if precedence(top) >= precedence(char) {
							postfix.WriteRune(stack.Remove(stack.Back()).(rune))
							postfix.WriteRune(' ')
						} else {
							break
						}
					}
					stack.PushBack(char)
					isUnary = true
				}
			}
		} else {
			return "", errors.New("invalid data")
		}
	}

	if num.Len() > 0 {
		postfix.WriteString(num.String() + " ")
	}

	for stack.Len() > 0 {
		top := stack.Remove(stack.Back()).(rune)
		postfix.WriteRune(top)
		postfix.WriteRune(' ')
	}

	return strings.TrimSpace(postfix.String()), nil
}
