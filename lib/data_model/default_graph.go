package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
)

// DefaultGraphValue is a constant that represents the default graph value.
const DefaultGraphValue = "rdfgo-DefaultGraph"

// DefaultGraphString is a constant that represents the default graph string.
const DefaultGraphString = "<>"

type defaultGraph struct{}

// NewDefaultGraph is a constructor for creating a default graph. This returns an implementation of the interfaces.IDefaultGraph interface.
func NewDefaultGraph() interfaces.IDefaultGraph {
	return &defaultGraph{}
}

// Equals is a method that checks if two default graphs are equal. This returns true if the other term is also a default graph.
func (d *defaultGraph) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	return interfaces.DefaultGraphType == other.GetType()
}

// GetType is a method that returns the (integer) type of the default graph.
func (d *defaultGraph) GetType() interfaces.TermType {
	return interfaces.DefaultGraphType
}

// GetValue is a method that returns the value of the default graph (DefaultGraphValue).
func (d *defaultGraph) GetValue() string {
	return DefaultGraphValue
}

// ToString is a method that returns the string representation of the default graph (DefaultGraphString).
func (d *defaultGraph) ToString() string {
	return DefaultGraphString
}
