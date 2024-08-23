package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
	"testing"
)

func TestLiteralCreation(t *testing.T) {
	tests := []struct {
		name             string
		literal          interfaces.ILiteral
		expectedValue    string
		expectedLang     string
		expectedDatatype interfaces.INamedNode
	}{
		{
			name:             "NewLiteral with language and datatype",
			literal:          NewLiteral("l1", "en", NewNamedNode("http://example.com")),
			expectedValue:    "l1",
			expectedLang:     "en",
			expectedDatatype: NewNamedNode("http://example.com"),
		},
		{
			name:             "NewStringLiteral with language",
			literal:          NewStringLiteral("l1", "en"),
			expectedValue:    "l1",
			expectedLang:     "en",
			expectedDatatype: IRI.XSD.String,
		},
		{
			name:             "NewIntegerLiteral",
			literal:          NewIntegerLiteral(123),
			expectedValue:    "123",
			expectedLang:     "",
			expectedDatatype: IRI.XSD.Integer,
		},
		{
			name:             "NewDecimalLiteral",
			literal:          NewDecimalLiteral(123.456),
			expectedValue:    "123.456",
			expectedLang:     "",
			expectedDatatype: IRI.XSD.Decimal,
		},
		{
			name:             "NewDoubleLiteral",
			literal:          NewDoubleLiteral(123.456),
			expectedValue:    "123.456",
			expectedLang:     "",
			expectedDatatype: IRI.XSD.Double,
		},
		{
			name:             "NewBooleanLiteral",
			literal:          NewBooleanLiteral(true),
			expectedValue:    "true",
			expectedLang:     "",
			expectedDatatype: IRI.XSD.Boolean,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.literal.GetValue() != tt.expectedValue {
				t.Errorf("Literal value should be %s, but got %s", tt.expectedValue, tt.literal.GetValue())
			}
			if tt.literal.GetLanguage() != tt.expectedLang {
				t.Errorf("Literal language should be %s, but got %s", tt.expectedLang, tt.literal.GetLanguage())
			}
			if !tt.literal.GetDatatype().Equals(tt.expectedDatatype) {
				t.Errorf(
					"Literal datatype should equal <%s>, but got %s",
					tt.expectedDatatype,
					tt.literal.GetDatatype().ToString(),
				)
			}
		})
	}
}

func TestLiteral_GetType(t *testing.T) {
	l1 := NewLiteral("l1", "en", NewNamedNode("http://example.com"))
	if l1.GetType() != interfaces.LiteralType {
		t.Errorf("Literal type should be %s", interfaces.LiteralType)
	}
}

func TestLiteral_GetValue(t *testing.T) {
	l1 := NewLiteral("l1", "en", NewNamedNode("http://example.com"))
	if l1.GetValue() != "l1" {
		t.Errorf("Literal value should be l1")
	}
}

func TestLiteral_GetDatatype(t *testing.T) {
	l1 := NewLiteral("l1", "en", NewNamedNode("http://example.com"))
	if !l1.GetDatatype().Equals(NewNamedNode("http://example.com")) {
		t.Errorf("Literal datatype should equal <http://example.com>")
	}
}

func TestLiteral_GetLanguage(t *testing.T) {
	l1 := NewLiteral("l1", "en", NewNamedNode("http://example.com"))
	if l1.GetLanguage() != "en" {
		t.Errorf("Literal language equal 'en'")
	}
}

func TestLiteral_Equals(t *testing.T) {
	l1 := NewLiteral("l1", "", nil)
	l2 := NewLiteral("l2", "", nil)
	l3 := NewLiteral("l1", "", nil)
	l4 := NewLiteral("l1", "en", nil)
	l5 := NewLiteral("l1", "", NewNamedNode("http://example.com"))
	l6 := NewLiteral("l1", "", NewNamedNode("http://example.com"))
	l7 := NewNamedNode("l1")
	if !l1.Equals(l1) {
		t.Errorf("Literal should equal itself")
	}
	if l1.Equals(l2) {
		t.Errorf("Literal should not equal another Literal")
	}
	if !l1.Equals(l3) {
		t.Errorf("Literal should equal another Literal with same value")
	}
	if l1.Equals(l4) {
		t.Errorf("Literal should not equal another Literal with different language")
	}
	if l1.Equals(l5) {
		t.Errorf("Literal should not equal another Literal with different datatype")
	}
	if !l5.Equals(l6) {
		t.Errorf("Literal should equal another Literal with same datatype")
	}
	if l5.Equals(l7) {
		t.Errorf("Literal should not equal NamedNodes with same value")
	}
}

func TestLiteral_EqualsNil(t *testing.T) {
	l1 := NewLiteral("l1", "", nil)
	if l1.Equals(nil) {
		t.Errorf("Literal should not equal nil")
	}
}

func TestLiteralToString(t *testing.T) {
	tests := []struct {
		name     string
		literal  Literal
		expected string
	}{
		{
			name: "Without language tag",
			literal: Literal{
				value:    "example",
				language: "",
				datatype: NewNamedNode("http://example.com/datatype"),
			},
			expected: "\"example\"^^<http://example.com/datatype>",
		},
		{
			name: "With language tag",
			literal: Literal{
				value:    "example",
				language: "en",
				datatype: NewNamedNode("http://example.com/datatype"),
			},
			expected: "\"example\"@en^^<http://example.com/datatype>",
		},
		{
			name: "With different datatype",
			literal: Literal{
				value:    "123",
				language: "",
				datatype: NewNamedNode("http://example.com/integer"),
			},
			expected: "\"123\"^^<http://example.com/integer>",
		},
		{
			name: "With empty value",
			literal: Literal{
				value:    "",
				language: "",
				datatype: nil,
			},
			expected: "\"\"",
		},
		{
			name: "With language tag and empty value",
			literal: Literal{
				value:    "",
				language: "es",
				datatype: nil,
			},
			expected: "\"\"@es",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.literal.ToString()
			if actual != tt.expected {
				t.Errorf("ToString() = %v, want %v", actual, tt.expected)
			}
		})
	}
}
