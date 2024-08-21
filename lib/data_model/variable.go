package rdfgo

import (
	"fmt"
	"rdfgo/interfaces"
)

type Variable struct {
	value string
}

func NewVariable(value string) interfaces.IVariable {
	if value[0] == '?' {
		value = value[1:]
	}
	return &Variable{
		value: value,
	}
}

func (v *Variable) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	return interfaces.VariableType == other.GetType() && v.value == other.GetValue()
}

func (v *Variable) GetType() interfaces.TermType {
	return interfaces.VariableType
}

func (v *Variable) GetValue() string {
	return v.value
}

func (v *Variable) ToString() string {
	return fmt.Sprintf("?%s", v.value)
}
