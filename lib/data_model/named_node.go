package rdfgo

import (
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
)

type namedNode struct {
	value string
}

// NewNamedNode is a constructor for creating a named node. This returns an implementation of the interfaces.INamedNode interface.
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

// Equals checks if the current named node is equal to another term.
func (n *namedNode) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	return interfaces.NamedNodeType == other.GetType() && n.value == other.GetValue()
}

// GetType is a method that returns the (integer) type of the named node.
func (n *namedNode) GetType() interfaces.TermType {
	return interfaces.NamedNodeType
}

// GetValue returns the value of the named node (http://example.org).
func (n *namedNode) GetValue() string {
	return n.value
}

// ToString returns the string representation of the named node (<http://example.org>).
func (n *namedNode) ToString() string {
	return fmt.Sprintf("<%s>", n.value)
}
