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

func infixToPostfix(infix string) (string, error) {
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
