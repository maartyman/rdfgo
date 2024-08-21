package rdfgo

import (
	"fmt"
	"rdfgo/interfaces"
)

var blankNodeCounter = 0

type BlankNode struct {
	value string
}

func NewBlankNode(value string) interfaces.IBlankNode {
	for {
		if len(value) == 0 || (value[0] != '_' && value[0] != ':') {
			break
		}
		value = value[1:]
	}
	if value == "" {
		value = fmt.Sprintf("n3-%d", blankNodeCounter)
		blankNodeCounter++
	}
	return &BlankNode{
		value: value,
	}
}

func (b *BlankNode) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	return interfaces.BlankNodeType == other.GetType() && b.value == other.GetValue()
}

func (b *BlankNode) GetType() interfaces.TermType {
	return interfaces.BlankNodeType
}

func (b *BlankNode) GetValue() string {
	return b.value
}

func (b *BlankNode) ToString() string {
	return fmt.Sprintf("_:%s", b.value)
}
