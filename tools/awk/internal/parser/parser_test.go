package parser_test

import (
	"goawk/internal/parser"
	"testing"
)

func TestParser_ParseProgram(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		validate func(t *testing.T, pa parser.PatternAction, err error)
	}{
		{
			name:  "simple field print",
			input: "$1",
			validate: func(t *testing.T, pa parser.PatternAction, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if !pa.HasPattern || !pa.HasAction {
					t.Errorf("Expected HasPattern and HasAction to be true")
				}
				if pa.PrintAll {
					t.Errorf("Expected PrintAll to be false")
				}
				if len(pa.PrintFields) != 1 || pa.PrintFields[0] != 1 {
					t.Errorf("Expected PrintFields to be [1], got %v", pa.PrintFields)
				}
			},
		},
		{
			name:  "condition with action",
			input: "$1 == \"test\" { print $2 }",
			validate: func(t *testing.T, pa parser.PatternAction, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if !pa.HasPattern || !pa.HasAction {
					t.Errorf("Expected HasPattern and HasAction to be true")
				}
				if pa.Pattern.Field1 != 1 {
					t.Errorf("Expected Field1 to be 1, got %d", pa.Pattern.Field1)
				}
				if pa.Pattern.CompareOp != "==" {
					t.Errorf("Expected CompareOp to be '==', got '%s'", pa.Pattern.CompareOp)
				}
				if pa.Pattern.String2 != "test" {
					t.Errorf("Expected String2 to be 'test', got '%s'", pa.Pattern.String2)
				}
			},
		},
		{
			name:  "numeric comparison",
			input: "$2 > 100",
			validate: func(t *testing.T, pa parser.PatternAction, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if pa.Pattern.Field1 != 2 {
					t.Errorf("Expected Field1 to be 2, got %d", pa.Pattern.Field1)
				}
				if pa.Pattern.CompareOp != ">" {
					t.Errorf("Expected CompareOp to be '>', got '%s'", pa.Pattern.CompareOp)
				}
				if pa.Pattern.String2 != "100" {
					t.Errorf("Expected String2 to be '100', got '%s'", pa.Pattern.String2)
				}
			},
		},
		{
			name:  "empty program",
			input: "",
			validate: func(t *testing.T, pa parser.PatternAction, err error) {
				if err == nil {
					t.Errorf("Expected error for empty program")
				}
			},
		},
		{
			name:  "invalid field in print args",
			input: "{ print $abc }",
			validate: func(t *testing.T, pa parser.PatternAction, err error) {
				if err == nil {
					t.Errorf("Expected error for invalid field")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := parser.NewParser()
			pa, err := p.ParseProgram(tc.input)
			tc.validate(t, pa, err)
		})
	}
}

func TestParser_ParsePrintArgs(t *testing.T) {
	p := parser.NewParser()
	_, err := p.ParseProgram("test")
	if err != nil {
		// This is just to get access to parsePrintArgs
	}

	// Note: parsePrintArgs is private, so we test it through ParseProgram
	testCases := []struct {
		name        string
		program     string
		expectError bool
	}{
		{"single field", "{ print $1 }", false},
		{"multiple fields", "{ print $1, $2 }", false},
		{"with raw string", "{ print \"hello\" }", false},
		{"invalid field", "{ print $abc }", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := parser.NewParser()
			_, err := p.ParseProgram(tc.program)
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
