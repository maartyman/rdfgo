package rdfgo

import (
	"rdfgo/interfaces"
	"testing"
)

func TestDefaultGraph(t *testing.T) {
	dg := NewDefaultGraph()
	if dg.GetType() != interfaces.DefaultGraphType {
		t.Errorf("DefaultGraph type should be %s", interfaces.DefaultGraphType)
	}
	if dg.GetValue() != DefaultGraphValue {
		t.Errorf("DefaultGraph name should be %s", DefaultGraphValue)
	}
}

func TestDefaultGraph_Equals(t *testing.T) {
	dg1 := NewDefaultGraph()
	dg2 := NewDefaultGraph()
	dg3 := NewNamedNode("http://example.com")
	dg4 := NewNamedNode("")
	dg5 := NewNamedNode("graph")
	if !dg1.Equals(dg1) {
		t.Errorf("DefaultGraph should equal itself")
	}
	if !dg1.Equals(dg2) {
		t.Errorf("DefaultGraph should equal another DefaultGraph")
	}
	if dg1.Equals(dg3) {
		t.Errorf("DefaultGraph should not equal a NamedNode")
	}
	if dg1.Equals(dg4) {
		t.Errorf("DefaultGraph should not equal a empty NamedNode")
	}
	if dg1.Equals(dg5) {
		t.Errorf("DefaultGraph should not equal a NamedNode with value 'graph'")
	}
}

func TestDefaultGraph_EqualsNil(t *testing.T) {
	dg1 := NewDefaultGraph()
	if dg1.Equals(nil) {
		t.Errorf("DefaultGraph should not equal nil")
	}
}

func TestDefaultGraph_ToString(t *testing.T) {
	dg1 := NewDefaultGraph()
	if dg1.ToString() != DefaultGraphValue {
		t.Errorf("DefaultGraph to string should equal an empty string")
	}
}
