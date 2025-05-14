package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
	"testing"
)

func TestDefaultGraph(t *testing.T) {
	dg := NewDefaultGraph()
	if dg.GetType() != interfaces.DefaultGraphType {
		t.Errorf("defaultGraph type should be %s", interfaces.DefaultGraphType)
	}
	if dg.GetValue() != DefaultGraphValue {
		t.Errorf("defaultGraph name should be %s", DefaultGraphValue)
	}
}

func TestDefaultGraph_Equals(t *testing.T) {
	dg1 := NewDefaultGraph()
	dg2 := NewDefaultGraph()
	dg3 := NewNamedNode("http://example.com")
	dg4 := NewNamedNode("")
	dg5 := NewNamedNode("graph")
	if !dg1.Equals(dg1) {
		t.Errorf("defaultGraph should equal itself")
	}
	if !dg1.Equals(dg2) {
		t.Errorf("defaultGraph should equal another defaultGraph")
	}
	if dg1.Equals(dg3) {
		t.Errorf("defaultGraph should not equal a namedNode")
	}
	if dg1.Equals(dg4) {
		t.Errorf("defaultGraph should not equal a empty namedNode")
	}
	if dg1.Equals(dg5) {
		t.Errorf("defaultGraph should not equal a namedNode with value 'graph'")
	}
}

func TestDefaultGraph_EqualsNil(t *testing.T) {
	dg1 := NewDefaultGraph()
	if dg1.Equals(nil) {
		t.Errorf("defaultGraph should not equal nil")
	}
}

func TestDefaultGraph_ToString(t *testing.T) {
	dg1 := NewDefaultGraph()
	if dg1.ToString() != DefaultGraphString {
		t.Errorf("defaultGraph to string should equal an empty string")
	}
}
