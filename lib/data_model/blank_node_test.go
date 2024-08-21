package rdfgo

import (
	"rdfgo/interfaces"
	"testing"
)

func TestBlankNode_Equals(t *testing.T) {
	b1 := NewBlankNode("b1")
	b2 := NewBlankNode("b2")
	b3 := NewBlankNode("b1")
	b4 := NewBlankNode("")
	b5 := NewBlankNode("")

	tests := []struct {
		b1       interfaces.IBlankNode
		b2       interfaces.IBlankNode
		expected bool
		message  string
	}{
		{b1, b1, true, "A BlankNode should equal itself"},
		{b1, b2, false, "A BlankNode should not equal another BlankNode with a different value"},
		{b1, b3, true, "A BlankNode should equal another BlankNode with the same value"},
		{b4, b5, false, "Two blank nodes with no values should not equal each other"},
	}

	for _, tt := range tests {
		if tt.b1.Equals(tt.b2) != tt.expected {
			t.Errorf(tt.message)
		}
	}
}

func TestBlankNode_EqualsNil(t *testing.T) {
	b1 := NewBlankNode("b1")
	if b1.Equals(nil) {
		t.Errorf("BlankNode should not equal nil")
	}
}

func TestBlankNode_GetType(t *testing.T) {
	b := NewBlankNode("b")
	if b.GetType() != interfaces.BlankNodeType {
		t.Errorf("BlankNode type should be %s", interfaces.BlankNodeType)
	}
}

func TestBlankNode_GetValue(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		message  string
	}{
		{"_:b", "b", "The '_:' should be removed"},
		{"_b", "b", "The '_' should be removed"},
		{"___b", "b", "Multiple `_` should be removed"},
		{"_:___b", "b", "The '_:' and multiple `_` should be removed"},
		{"__:___b", "b", "The '__:' and multiple `_` should be removed"},
		{"b", "b", "No modification needed"},
	}

	for _, tt := range tests {
		b := NewBlankNode(tt.input)
		if b.GetValue() != tt.expected {
			t.Errorf("%s, but got %s", tt.message, b.GetValue())
		}
	}
}

func TestBlankNode_ToString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		message  string
	}{
		{"_:b", "_:b", "The '_:' should be removed"},
		{"_b", "_:b", "The '_' should be removed"},
		{"___b", "_:b", "Multiple `_` should be removed"},
		{"_:___b", "_:b", "The '_:' and multiple `_` should be removed"},
		{"__:___b", "_:b", "The '__:' and multiple `_` should be removed"},
		{"b", "_:b", "No modification needed"},
	}

	for _, tt := range tests {
		b := NewBlankNode(tt.input)
		if b.ToString() != tt.expected {
			t.Errorf("%s, but got %s", tt.message, b.ToString())
		}
	}
}
