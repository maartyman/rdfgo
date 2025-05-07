package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
)

const (
	xsd = "http://www.w3.org/2001/XMLSchema#"
	rdf = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
)

type xsdTerms struct {
	Decimal interfaces.INamedNode
	Boolean interfaces.INamedNode
	Double  interfaces.INamedNode
	Integer interfaces.INamedNode
	String  interfaces.INamedNode
}

type rdfTerms struct {
	Type       interfaces.INamedNode
	Nil        interfaces.INamedNode
	First      interfaces.INamedNode
	Rest       interfaces.INamedNode
	LangString interfaces.INamedNode
}

type terms struct {
	XSD xsdTerms
	RDF rdfTerms
}

var IRI = terms{
	XSD: xsdTerms{
		Decimal: NewNamedNode(xsd + "decimal"),
		Boolean: NewNamedNode(xsd + "boolean"),
		Double:  NewNamedNode(xsd + "double"),
		Integer: NewNamedNode(xsd + "integer"),
		String:  NewNamedNode(xsd + "string"),
	},
	RDF: rdfTerms{
		Type:       NewNamedNode(rdf + "type"),
		Nil:        NewNamedNode(rdf + "nil"),
		First:      NewNamedNode(rdf + "first"),
		Rest:       NewNamedNode(rdf + "rest"),
		LangString: NewNamedNode(rdf + "langString"),
	},
}
