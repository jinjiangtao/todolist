package lexer

import "strings"

type TokenType int

const (
	TokenField TokenType = iota
	TokenString
	TokenNumber
	TokenCompare
	TokenEqual
	TokenNotEqual
	TokenLT
	TokenGT
	TokenLTE
	TokenGTE
	TokenAnd
	TokenOr
	TokenNot
	TokenPrint
	TokenLBrace
	TokenRBrace
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF}
	}

	ch := l.input[l.pos]

	if ch == '$' && l.pos+1 < len(l.input) && l.input[l.pos+1] >= '0' && l.input[l.pos+1] <= '9' {
		l.pos++
		fieldNum := ""
		for l.pos < len(l.input) && l.input[l.pos] >= '0' && l.input[l.pos] <= '9' {
			fieldNum += string(l.input[l.pos])
			l.pos++
		}
		return Token{Type: TokenField, Value: fieldNum}
	}

	if ch == '"' {
		l.pos++
		str := ""
		for l.pos < len(l.input) && l.input[l.pos] != '"' {
			str += string(l.input[l.pos])
			l.pos++
		}
		if l.pos < len(l.input) {
			l.pos++
		}
		return Token{Type: TokenString, Value: str}
	}

	if ch == '\'' {
		l.pos++
		str := ""
		for l.pos < len(l.input) && l.input[l.pos] != '\'' {
			str += string(l.input[l.pos])
			l.pos++
		}
		if l.pos < len(l.input) {
			l.pos++
		}
		return Token{Type: TokenString, Value: str}
	}

	if l.pos+1 < len(l.input) {
		twoChar := string(ch) + string(l.input[l.pos+1])
		switch twoChar {
		case "==":
			l.pos += 2
			return Token{Type: TokenEqual, Value: "=="}
		case "!=":
			l.pos += 2
			return Token{Type: TokenNotEqual, Value: "!="}
		case "<=":
			l.pos += 2
			return Token{Type: TokenLTE, Value: "<="}
		case ">=":
			l.pos += 2
			return Token{Type: TokenGTE, Value: ">="}
		case "&&":
			l.pos += 2
			return Token{Type: TokenAnd, Value: "&&"}
		case "||":
			l.pos += 2
			return Token{Type: TokenOr, Value: "||"}
		}
	}

	switch ch {
	case '<':
		l.pos++
		return Token{Type: TokenLT, Value: "<"}
	case '>':
		l.pos++
		return Token{Type: TokenGT, Value: ">"}
	case '{':
		l.pos++
		return Token{Type: TokenLBrace, Value: "{"}
	case '}':
		l.pos++
		return Token{Type: TokenRBrace, Value: "}"}
	}

	if (ch >= '0' && ch <= '9') || ch == '.' {
		num := ""
		for l.pos < len(l.input) && ((l.input[l.pos] >= '0' && l.input[l.pos] <= '9') || l.input[l.pos] == '.') {
			num += string(l.input[l.pos])
			l.pos++
		}
		return Token{Type: TokenNumber, Value: num}
	}

	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
		ident := ""
		for l.pos < len(l.input) && ((l.input[l.pos] >= 'a' && l.input[l.pos] <= 'z') ||
			(l.input[l.pos] >= 'A' && l.input[l.pos] <= 'Z') ||
			(l.input[l.pos] >= '0' && l.input[l.pos] <= '9') ||
			l.input[l.pos] == '_') {
			ident += string(l.input[l.pos])
			l.pos++
		}
		if strings.ToLower(ident) == "print" {
			return Token{Type: TokenPrint, Value: ident}
		}
		return Token{Type: TokenString, Value: ident}
	}

	l.pos++
	return l.NextToken()
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && (l.input[l.pos] == ' ' || l.input[l.pos] == '\t') {
		l.pos++
	}
}
