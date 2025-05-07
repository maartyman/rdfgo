package rdfgo

import (
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
)

type variable struct {
	value string
}

func NewVariable(value string) interfaces.IVariable {
	if value[0] == '?' {
		value = value[1:]
	}
	return &variable{
		value: value,
	}
}

func (v *variable) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	return interfaces.VariableType == other.GetType() && v.value == other.GetValue()
}

func (v *variable) GetType() interfaces.TermType {
	return interfaces.VariableType
}

func (v *variable) GetValue() string {
	return v.value
}

func (v *variable) ToString() string {
	return fmt.Sprintf("?%s", v.value)
}
