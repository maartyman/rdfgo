package rdfgo

import (
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
)

type variable struct {
	value string
}

// NewVariable is a constructor for creating a variable. This returns an implementation of the interfaces.IVariable interface.
func NewVariable(value string) interfaces.IVariable {
	if value[0] == '?' {
		value = value[1:]
	}
	return &variable{
		value: value,
	}
}

// Equals checks if the current variable is equal to another term.
func (v *variable) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	return interfaces.VariableType == other.GetType() && v.value == other.GetValue()
}

// GetType is a method that returns the (integer) type of the variable.
func (v *variable) GetType() interfaces.TermType {
	return interfaces.VariableType
}

// GetValue returns the value of the variable (v1).
func (v *variable) GetValue() string {
	return v.value
}

// ToString returns the string representation of the variable (?v1).
func (v *variable) ToString() string {
	return fmt.Sprintf("?%s", v.value)
}
