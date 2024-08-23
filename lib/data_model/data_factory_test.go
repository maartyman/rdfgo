package rdfgo

import (
	"testing"
)

func TestDataFactory_NamedNode(t *testing.T) {
	df := NewDataFactory()
	if df.NamedNode("http://example.org").Equals(NewNamedNode("http://example.org")) == false {
		t.Error("Data factory NamedNode should return a NamedNode with the given value")
	}
}

func TestDataFactory_BlankNode(t *testing.T) {
	df := NewDataFactory()
	if df.BlankNode("b1").Equals(NewBlankNode("b1")) == false {
		t.Error("Data factory BlankNode should return a BlankNode with the given value")
	}
	if df.BlankNode("").Equals(NewBlankNode("")) {
		t.Error("A blank node build by df and a blank node build by NewBlankNode should have different counters")
	}
	if df.BlankNode("").Equals(df.BlankNode("")) {
		t.Error("When no value is passed it should return different blank nodes")
	}
}

func TestDataFactory_SimpleLiteral(t *testing.T) {
	df := NewDataFactory()
	if !df.SimpleLiteral("value").
		Equals(NewLiteral("value", "", df.NamedNode("http://www.w3.org/2001/XMLSchema#string"))) {
		t.Error("Data factory SimpleLiteral should return a Literal with the given value and default datatype")
	}
}

func TestDataFactory_Literal(t *testing.T) {
	df := NewDataFactory()
	if !df.Literal("value", "en", nil).
		Equals(NewLiteral("value", "en", df.NamedNode("http://www.w3.org/2001/XMLSchema#string"))) {
		t.Error("Data factory Literal should return a Literal with the given value, language and default datatype")
	}
	if !df.Literal("value", "en", df.NamedNode("http://www.w3.org/2001/XMLSchema#string")).
		Equals(NewLiteral("value", "en", df.NamedNode("http://www.w3.org/2001/XMLSchema#string"))) {
		t.Error("Data factory Literal should return a Literal with the given value, language and datatype")
	}
}

func TestDataFactory_Variable(t *testing.T) {
	df := NewDataFactory()
	if df.Variable("v").Equals(NewVariable("v")) == false {
		t.Error("Data factory Variable should return a Variable with the given value")
	}
}

func TestDataFactory_DefaultGraph(t *testing.T) {
	df := NewDataFactory()
	if df.DefaultGraph().Equals(NewDefaultGraph()) == false {
		t.Error("Data factory DefaultGraph should return a DefaultGraph")
	}
}

func TestDataFactory_Quad(t *testing.T) {
	df := NewDataFactory()
	dfQuad, _ := df.Quad(df.NamedNode("s"), df.NamedNode("p"), df.NamedNode("o"), df.DefaultGraph())
	originalQuad, _ := NewQuad(df.NamedNode("s"), df.NamedNode("p"), df.NamedNode("o"), df.DefaultGraph())
	if dfQuad.Equals(originalQuad) == false {
		t.Error("Data factory Quad should return a Quad with the given terms")
	}
}
