package lexer_test

import (
	"goawk/internal/lexer"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{
			name:  "single field",
			input: "$1",
			expected: []lexer.Token{
				{Type: lexer.TokenField, Value: "1"},
				{Type: lexer.TokenEOF},
			},
		},
		{
			name:  "multiple fields",
			input: "$1 $2",
			expected: []lexer.Token{
				{Type: lexer.TokenField, Value: "1"},
				{Type: lexer.TokenField, Value: "2"},
				{Type: lexer.TokenEOF},
			},
		},
		{
			name:  "string literal",
			input: "\"hello\"",
			expected: []lexer.Token{
				{Type: lexer.TokenString, Value: "hello"},
				{Type: lexer.TokenEOF},
			},
		},
		{
			name:  "single quote string",
			input: "'world'",
			expected: []lexer.Token{
				{Type: lexer.TokenString, Value: "world"},
				{Type: lexer.TokenEOF},
			},
		},
		{
			name:  "numeric value",
			input: "123.45",
			expected: []lexer.Token{
				{Type: lexer.TokenNumber, Value: "123.45"},
				{Type: lexer.TokenEOF},
			},
		},
		{
			name:  "print keyword",
			input: "print",
			expected: []lexer.Token{
				{Type: lexer.TokenPrint, Value: "print"},
				{Type: lexer.TokenEOF},
			},
		},
		{
			name:  "comparison operators",
			input: "== != <= >= < >",
			expected: []lexer.Token{
				{Type: lexer.TokenEqual, Value: "=="},
				{Type: lexer.TokenNotEqual, Value: "!="},
				{Type: lexer.TokenLTE, Value: "<="},
				{Type: lexer.TokenGTE, Value: ">="},
				{Type: lexer.TokenLT, Value: "<"},
				{Type: lexer.TokenGT, Value: ">"},
				{Type: lexer.TokenEOF},
			},
		},
		{
			name:  "logical operators",
			input: "&& ||",
			expected: []lexer.Token{
				{Type: lexer.TokenAnd, Value: "&&"},
				{Type: lexer.TokenOr, Value: "||"},
				{Type: lexer.TokenEOF},
			},
		},
		{
			name:  "braces",
			input: "{}",
			expected: []lexer.Token{
				{Type: lexer.TokenLBrace, Value: "{"},
				{Type: lexer.TokenRBrace, Value: "}"},
				{Type: lexer.TokenEOF},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := lexer.NewLexer(tc.input)
			for _, exp := range tc.expected {
				token := l.NextToken()
				if token.Type != exp.Type {
					t.Errorf("Expected token type %v, got %v (value: '%s')", exp.Type, token.Type, token.Value)
				}
				if token.Value != exp.Value {
					t.Errorf("Expected token value '%s', got '%s'", exp.Value, token.Value)
				}
			}
		})
	}
}

func TestLexer_SkipWhitespace(t *testing.T) {
	input := "   $1  print  \"test\"   "
	expected := []lexer.Token{
		{Type: lexer.TokenField, Value: "1"},
		{Type: lexer.TokenPrint, Value: "print"},
		{Type: lexer.TokenString, Value: "test"},
		{Type: lexer.TokenEOF},
	}

	l := lexer.NewLexer(input)
	for _, exp := range expected {
		token := l.NextToken()
		if token.Type != exp.Type || token.Value != exp.Value {
			t.Errorf("Expected %v (value: '%s'), got %v (value: '%s')", exp.Type, exp.Value, token.Type, token.Value)
		}
	}
}
