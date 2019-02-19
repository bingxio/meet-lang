package lexer

import (
	"meet-lang/src/token"
)

type Lexer struct {
	input   string
	char    string
	current int
	length  int
	tokens  []token.Token
}

func New(input string) []token.Token {
	l := Lexer{
		input:   input,
		current: 0,
	}
	l.tokenizer()

	return l.tokens
}

func (l *Lexer) tokenizer() {
	l.length = len([]rune(l.input))

	for l.current < l.length {
		l.char = string([]rune(l.input)[l.current])

		if l.char == " " || l.char == "\t" || l.char == "\n" || l.char == "\r" {
			l.refreshNextChar()
			continue
		}

		if l.char == "#" {
			l.refreshNextChar()

			for l.char != "#" {
				if !l.refreshNextChar() {
					panic("注释必须以 # 结束，并一一对应")
				}
			}

			l.refreshNextChar()
			continue
		}

		////////////////////////////////////

		if isString(l.char) {
			value := ""

			l.refreshNextChar()

			for !isString(l.char) {
				value += l.char

				if !l.refreshNextChar() {
					panic("字符串必须一一对应")
				}
			}

			l.newToken(token.STRING, value)
			l.refreshNextChar()

			continue
		}

		if isLetter(l.char) {
			value := ""

			for isLetter(l.char) {
				value += l.char

				if !l.refreshNextChar() {
					break
				}
			}

			if l.char == "[" {
				list := ""

				for l.char != "]" {
					list += l.char

					if !l.refreshNextChar() {
						break
					}
				}

				list += "]"
				value += list

				l.refreshNextChar()
				l.newToken(token.LIST, value)

				continue
			}

			l.newToken(token.NAME, value)

			continue
		}

		if isDigit(l.char) {
			value := ""

			for isDigit(l.char) {
				value += l.char

				if !l.refreshNextChar() {
					break
				}
			}

			l.newToken(token.DIGIT, value)

			continue
		}

		////////////////////////////////////

		if l.char == "-" {
			if !l.refreshNextChar() {
				panic("判断运算符出错，可能是其他类型，而缺少了一些符号：" + l.char)
			}

			if l.char == ">" {
				l.newToken(token.POINTER, "->")
				l.refreshNextChar()

				continue
			}

			if l.char == "=" {
				l.newToken(token.MINUS_ASSIGN, "-=")
				l.refreshNextChar()

				continue
			}

			if l.char == "-" {
				l.newToken(token.MINUS_ONE, "--")
				l.refreshNextChar()

				continue
			}

			l.newToken(token.MINUS, "-")
			l.refreshNextChar()

			continue
		}

		if l.char == "=" {
			if !l.refreshNextChar() {
				panic("判断运算符出错，可能是其他类型，而缺少了一些符号：" + l.char)
			}

			if l.char == ">" {
				l.newToken(token.FUNCTION_POINTER, "=>")
				l.refreshNextChar()

				continue
			}

			if l.char == "=" {
				l.newToken(token.EQ, "==")
				l.refreshNextChar()

				continue
			}

			l.newToken(token.ASSIGN, "=")
			l.refreshNextChar()

			continue
		}

		if l.char == "!" {
			if !l.refreshNextChar() {
				panic("判断运算符出错，可能是其他类型，而缺少了一些符号：" + l.char)
			}

			if l.char == "=" {
				l.newToken(token.NOT_EQ, "!=")
				l.refreshNextChar()

				continue
			}

			l.newToken(token.BANG, "!")
			l.refreshNextChar()

			continue
		}

		if l.char == ">" {
			if !l.refreshNextChar() {
				panic("判断运算符出错，可能是其他类型，而缺少了一些符号：" + l.char)
			}

			if l.char == "=" {
				l.newToken(token.LT_ASSIGN, ">=")
				l.refreshNextChar()

				continue
			}

			l.newToken(token.LT, ">")
			l.refreshNextChar()

			continue
		}

		if l.char == "<" {
			if !l.refreshNextChar() {
				panic("判断运算符出错，可能是其他类型，而缺少了一些符号：" + l.char)
			}

			if l.char == "=" {
				l.newToken(token.RT_ASSIGN, "<=")
				l.refreshNextChar()

				continue
			}

			l.newToken(token.RT, "<")
			l.refreshNextChar()

			continue
		}

		if l.char == "+" {
			if !l.refreshNextChar() {
				panic("判断运算符出错，可能是其他类型，而缺少了一些符号：" + l.char)
			}

			if l.char == "=" {
				l.newToken(token.PLUS_ASSIGN, "+=")
				l.refreshNextChar()

				continue
			}

			if l.char == "+" {
				l.newToken(token.PLUS_ONE, "++")
				l.refreshNextChar()

				continue
			}

			l.newToken(token.PLUS, "+")
			l.refreshNextChar()

			continue
		}

		if l.char == "*" {
			if !l.refreshNextChar() {
				panic("判断运算符出错，可能是其他类型，而缺少了一些符号：" + l.char)
			}

			if l.char == "=" {
				l.newToken(token.ASTERISK_ASSIGN, "*=")
				l.refreshNextChar()

				continue
			}

			l.newToken(token.ASTERISK, "*")
			l.refreshNextChar()

			continue
		}

		if l.char == "/" {
			if !l.refreshNextChar() {
				panic("判断运算符出错，可能是其他类型，而缺少了一些符号：" + l.char)
			}

			if l.char == "=" {
				l.newToken(token.DIV_ASSIGN, "/=")
				l.refreshNextChar()

				continue
			}

			l.newToken(token.DIV, "/")
			l.refreshNextChar()

			continue
		}

		if l.char == "(" || l.char == ")" || l.char == "{" || l.char == "}" || l.char == ";" ||
			l.char == "," || l.char == "%" || l.char == "[" || l.char == "]" {
			value := l.char

			switch value {
			case "(":
				l.newToken(token.LPAREN, value)
			case ")":
				l.newToken(token.RPAREN, value)
			case "{":
				l.newToken(token.LBRACE, value)
			case "}":
				l.newToken(token.RBRACE, value)
			case ";":
				l.newToken(token.SEMICOLON, value)
			case ",":
				l.newToken(token.COMMA, value)
			case "%":
				l.newToken(token.MODULAR, value)
			case "[":
				l.newToken(token.LBRACKET, value)
			case "]":
				l.newToken(token.RBRACKET, value)
			}
			l.refreshNextChar()

			continue
		}

		panic("我不知道这个字符是什么：" + l.char)
	}

	l.newToken(token.EOF, "EOF")
}

func (l *Lexer) refreshNextChar() bool {
	l.current++

	if l.current == l.length {
		return false
	}

	l.char = string([]rune(l.input)[l.current])

	return true
}

func (l *Lexer) newToken(types, value string) {
	l.tokens = append(l.tokens, token.Token{
		Type:  types,
		Value: value,
	})
}

func isLetter(char string) bool {
	ch := []byte(char)[0]
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(char string) bool {
	ch := []byte(char)[0]
	return '0' <= ch && ch <= '9'
}

func isString(char string) bool {
	return char == "'"
}
