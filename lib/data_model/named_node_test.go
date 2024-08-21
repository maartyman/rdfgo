package rdfgo

import (
	"rdfgo/interfaces"
	"testing"
)

func TestNamedNode_NewNamedNode(t *testing.T) {
	l1 := NewNamedNode("<l1>")
	l2 := NewNamedNode("l1")
	if l1.GetValue() != "l1" {
		t.Errorf("NamedNode value should be l1")
	}
	if !l1.Equals(l2) {
		t.Errorf("NamedNode with and without `<` and `>` should be equal")
	}
}

func TestNamedNode_GetType(t *testing.T) {
	l1 := NewNamedNode("l1")
	if l1.GetType() != interfaces.NamedNodeType {
		t.Errorf("NamedNode type should be %s", interfaces.NamedNodeType)
	}
}

func TestNamedNode_GetValue(t *testing.T) {
	l1 := NewNamedNode("l1")
	if l1.GetValue() != "l1" {
		t.Errorf("NamedNode value should be l1")
	}
}

func TestNamedNode_Equals(t *testing.T) {
	l1 := NewNamedNode("l1")
	l2 := NewNamedNode("l2")
	l3 := NewNamedNode("l1")
	if !l1.Equals(l1) {
		t.Errorf("l1 should not equal l1")
	}
	if l1.Equals(l2) {
		t.Errorf("l1 should not equal l2")
	}
	if !l1.Equals(l3) {
		t.Errorf("l1 should equal l3")
	}
}

func TestNamedNode_EqualsNil(t *testing.T) {
	l1 := NewNamedNode("l1")
	if l1.Equals(nil) {
		t.Errorf("NameNode should not equal nil")
	}
}

func TestNamedNode_ToString(t *testing.T) {
	l1 := NewNamedNode("l1")
	if l1.ToString() != "<l1>" {
		t.Errorf("NameNode to string should equal <l1>")
	}
}
