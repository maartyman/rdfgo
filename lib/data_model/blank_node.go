package rdfgo

import (
	"fmt"
	"github.com/maartyman/rdfgo/interfaces"
)

var blankNodeCounter = 0

type blankNode struct {
	value string
}

// NewBlankNode is a constructor for creating a blank node. This returns an implementation of the interfaces.IBlankNode interface.
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
	return &blankNode{
		value: value,
	}
}

// Equals checks if the current blank node is equal to another term.
func (b *blankNode) Equals(other interfaces.ITerm) bool {
	if other == nil {
		return false
	}
	return interfaces.BlankNodeType == other.GetType() && b.value == other.GetValue()
}

// GetType is a method that returns the (integer) type of the blank node.
func (b *blankNode) GetType() interfaces.TermType {
	return interfaces.BlankNodeType
}

// GetValue returns the value of the blank node (b1).
func (b *blankNode) GetValue() string {
	return b.value
}

// ToString returns the string representation of the blank node (_:b1).
func (b *blankNode) ToString() string {
	return fmt.Sprintf("_:%s", b.value)
}
