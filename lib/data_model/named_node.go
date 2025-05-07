package rdfgo

import (
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
)

type namedNode struct {
	value string
}

func NewNamedNode(value string) interfaces.INamedNode {
	if len(value) > 0 {
		if value[0] == '<' {
			value = value[1:]
		}
	}
	if len(value) > 0 {
		if value[len(value)-1] == '>' {
			value = value[:len(value)-1]
		}
	}
	return &namedNode{
		value: value,
	}
}

func (n *namedNode) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	return interfaces.NamedNodeType == other.GetType() && n.value == other.GetValue()
}

func (n *namedNode) GetType() interfaces.TermType {
	return interfaces.NamedNodeType
}

func (n *namedNode) GetValue() string {
	return n.value
}

func (n *namedNode) ToString() string {
	return fmt.Sprintf("<%s>", n.value)
}
