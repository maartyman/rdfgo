package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
)

const DefaultGraphValue = "rdfgo-DefaultGraph"
const DefaultGraphString = "<>"

type DefaultGraph struct{}

func NewDefaultGraph() interfaces.IDefaultGraph {
	return &DefaultGraph{}
}

func (d *DefaultGraph) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	return interfaces.DefaultGraphType == other.GetType()
}

func (d *DefaultGraph) GetType() interfaces.TermType {
	return interfaces.DefaultGraphType
}

func (d *DefaultGraph) GetValue() string {
	return DefaultGraphValue
}

func (d *DefaultGraph) ToString() string {
	return DefaultGraphString
}
