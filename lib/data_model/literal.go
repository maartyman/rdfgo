package rdfgo

import (
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
)

type Literal struct {
	value    string
	language string
	datatype interfaces.INamedNode
}

func NewLiteral(value string, language string, datatype interfaces.INamedNode) interfaces.ILiteral {
	return &Literal{
		value:    value,
		language: language,
		datatype: datatype,
	}
}

func NewStringLiteral(value string, language string) interfaces.ILiteral {
	return NewLiteral(value, language, IRI.XSD.String)
}

func NewIntegerLiteral(value int) interfaces.ILiteral {
	return NewLiteral(fmt.Sprintf("%d", value), "", IRI.XSD.Integer)
}

func NewDecimalLiteral(value float64) interfaces.ILiteral {
	return NewLiteral(fmt.Sprintf("%g", value), "", IRI.XSD.Decimal)
}

func NewDoubleLiteral(value float64) interfaces.ILiteral {
	return NewLiteral(fmt.Sprintf("%g", value), "", IRI.XSD.Double)
}

func NewBooleanLiteral(value bool) interfaces.ILiteral {
	return NewLiteral(fmt.Sprintf("%t", value), "", IRI.XSD.Boolean)
}

func (l *Literal) Equals(other interfaces.ITerm) bool {
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

func (l *Literal) GetValue() string {
	return l.value
}

func (l *Literal) GetType() interfaces.TermType {
	return interfaces.LiteralType
}

func (l *Literal) GetLanguage() string {
	return l.language
}

func (l *Literal) GetDatatype() interfaces.INamedNode {
	return l.datatype
}

func (l *Literal) ToString() string {
	languageString := ""
	if l.language != "" {
		languageString = fmt.Sprintf("@%s", l.language)
	}
	dataTypeString := ""
	if l.datatype != nil {
		dataTypeString = fmt.Sprintf("^^%s", l.datatype.ToString())

	}
	return fmt.Sprintf("\"%s\"%s%s", l.value, languageString, dataTypeString)
}
