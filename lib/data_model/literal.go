package rdfgo

import (
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
	"strings"
)

type literal struct {
	value    string
	language string
	datatype interfaces.INamedNode
}

// NewLiteral is a constructor for creating a literal. This returns an implementation of the interfaces.ILiteral interface.
func NewLiteral(value string, language string, datatype interfaces.INamedNode) interfaces.ILiteral {
	if language != "" && datatype != nil {
		datatype = nil
	}
	return &literal{
		value:    value,
		language: language,
		datatype: datatype,
	}
}

// NewStringLiteral is a constructor for creating a string literal. This returns an implementation of the interfaces.ILiteral interface.
func NewStringLiteral(value string, language string) interfaces.ILiteral {
	return NewLiteral(value, language, IRI.XSD.String)
}

// NewIntegerLiteral is a constructor for creating an integer literal. This returns an implementation of the interfaces.ILiteral interface.
func NewIntegerLiteral(value int) interfaces.ILiteral {
	return NewLiteral(fmt.Sprintf("%d", value), "", IRI.XSD.Integer)
}

// NewDecimalLiteral is a constructor for creating a decimal literal. This returns an implementation of the interfaces.ILiteral interface.
func NewDecimalLiteral(value float64) interfaces.ILiteral {
	return NewLiteral(fmt.Sprintf("%g", value), "", IRI.XSD.Decimal)
}

// NewDoubleLiteral is a constructor for creating a double literal. This returns an implementation of the interfaces.ILiteral interface.
func NewDoubleLiteral(value float64) interfaces.ILiteral {
	return NewLiteral(fmt.Sprintf("%g", value), "", IRI.XSD.Double)
}

// NewBooleanLiteral is a constructor for creating a boolean literal. This returns an implementation of the interfaces.ILiteral interface.
func NewBooleanLiteral(value bool) interfaces.ILiteral {
	return NewLiteral(fmt.Sprintf("%t", value), "", IRI.XSD.Boolean)
}

// Equals is a method that checks if two literals are equal. It compares the value, language, and datatype of the literals.
func (l *literal) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	if l == other {
		return true
	}
	literal, ok := other.(interfaces.ILiteral)
	if !ok || interfaces.LiteralType != other.GetType() {
		return false
	}
	return l.value == literal.GetValue() && l.language == literal.GetLanguage() &&
		((l.datatype == nil && literal.GetDatatype() == nil) ||
			(l.datatype != nil && l.datatype.Equals(literal.GetDatatype())))
}

// GetValue is a method that returns the value of the literal ("test").
func (l *literal) GetValue() string {
	return l.value
}

// GetType is a method that returns the (integer) type of the literal.
func (l *literal) GetType() interfaces.TermType {
	return interfaces.LiteralType
}

// GetLanguage is a method that returns the language of the literal (en).
func (l *literal) GetLanguage() string {
	return l.language
}

// GetDatatype is a method that returns the datatype of the literal (xsd:string).
func (l *literal) GetDatatype() interfaces.INamedNode {
	return l.datatype
}

// ToString is a method that returns the string representation of the literal ("test" | "test"@en | "1"^^<xsd:integer>).
func (l *literal) ToString() string {
	escapedValue := strings.NewReplacer(
		`"`, `\"`,
		`\`, `\\`,
		"\n", `\n`,
		"\t", `\t`,
		"\r", `\r`,
	).Replace(l.value)

	if l.language != "" {
		return fmt.Sprintf("\"%s\"@%s", escapedValue, l.language)
	}
	if l.datatype == nil || l.datatype.Equals(IRI.XSD.String) {
		return fmt.Sprintf("\"%s\"", escapedValue)
	}
	return fmt.Sprintf("\"%s\"^^%s", escapedValue, l.datatype.ToString())
}
