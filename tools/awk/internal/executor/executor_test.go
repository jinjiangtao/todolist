package executor_test

import (
	"bytes"
	"goawk/internal/executor"
	"goawk/internal/parser"
	"testing"
)

func TestExecutor_Execute(t *testing.T) {
	testCases := []struct {
		name     string
		program  string
		input    string
		expected string
	}{
		{
			name:     "print first field",
			program:  "$1",
			input:    "John 25\nAlice 30",
			expected: "John\nAlice\n",
		},
		{
			name:     "print all",
			program:  "{ print }",
			input:    "hello\nworld",
			expected: "hello\nworld\n",
		},
		{
			name:     "numeric condition",
			program:  "$2 > 25 { print $1 }",
			input:    "John 25\nAlice 30\nBob 20",
			expected: "Alice\n",
		},
		{
			name:     "string condition",
			program:  "$1 == \"Alice\" { print $2 }",
			input:    "John 25\nAlice 30",
			expected: "30\n",
		},
		{
			name:     "chinese text",
			program:  "$1",
			input:    "你好 世界\n张三 25",
			expected: "你好\n张三\n",
		},
		{
			name:     "comma separator",
			program:  "$2",
			input:    "John,25,Engineer\nAlice,30,Manager",
			expected: "25\n30\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := parser.NewParser()
			pa, err := p.ParseProgram(tc.program)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			fs := " "
			if tc.name == "comma separator" {
				fs = ","
			}
			exe := executor.NewAWKExecutor(fs)
			exe.SetProgram(pa)

			input := bytes.NewBufferString(tc.input)
			var output bytes.Buffer

			if err := exe.Execute(input, &output); err != nil {
				t.Fatalf("Execute error: %v", err)
			}

			if output.String() != tc.expected {
				t.Errorf("Expected '%q', got '%q'", tc.expected, output.String())
			}
		})
	}
}
