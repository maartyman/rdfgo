package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
)

const DefaultGraphValue = "rdfgo-DefaultGraph"
const DefaultGraphString = "<>"

type defaultGraph struct{}

func NewDefaultGraph() interfaces.IDefaultGraph {
	return &defaultGraph{}
}

func (d *defaultGraph) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	return interfaces.DefaultGraphType == other.GetType()
}

func (d *defaultGraph) GetType() interfaces.TermType {
	return interfaces.DefaultGraphType
}

func (d *defaultGraph) GetValue() string {
	return DefaultGraphValue
}

func (d *defaultGraph) ToString() string {
	return DefaultGraphString
}
