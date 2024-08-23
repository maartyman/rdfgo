package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
)

const (
	xsd = "http://www.w3.org/2001/XMLSchema#"
	//rdf  = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	//swap = "http://www.w3.org/2000/10/swap/"
)

type XSDTerms struct {
	Decimal interfaces.INamedNode
	Boolean interfaces.INamedNode
	Double  interfaces.INamedNode
	Integer interfaces.INamedNode
	String  interfaces.INamedNode
}

/*
type RDFTerms struct {
	Type       interfaces.INamedNode
	Nil        interfaces.INamedNode
	First      interfaces.INamedNode
	Rest       interfaces.INamedNode
	LangString interfaces.INamedNode
}

type OWLTerms struct {
	SameAs interfaces.INamedNode
}

type RTerms struct {
	ForSome interfaces.INamedNode
	ForAll  interfaces.INamedNode
}

type LogTerms struct {
	Implies interfaces.INamedNode
}
*/

type Terms struct {
	XSD XSDTerms
	/*
		RDF RDFTerms
		OWL OWLTerms
		R   RTerms
		Log LogTerms
	*/
}

var IRI = Terms{
	XSD: XSDTerms{
		Decimal: NewNamedNode(xsd + "decimal"),
		Boolean: NewNamedNode(xsd + "boolean"),
		Double:  NewNamedNode(xsd + "double"),
		Integer: NewNamedNode(xsd + "integer"),
		String:  NewNamedNode(xsd + "string"),
	},
	/*
		RDF: RDFTerms{
			Type:       NewNamedNode(rdf + "type"),
			Nil:        NewNamedNode(rdf + "nil"),
			First:      NewNamedNode(rdf + "first"),
			Rest:       NewNamedNode(rdf + "rest"),
			LangString: NewNamedNode(rdf + "langString"),
		},
		OWL: OWLTerms{
			SameAs: NewNamedNode("http://www.w3.org/2002/07/owl#sameAs"),
		},
		R: RTerms{
			ForSome: NewNamedNode(swap + "reify#forSome"),
			ForAll:  swap + "reify#forAll"),
		},
		Log: LogTerms{
			Implies: NewNamedNode(swap + "log#implies"),
		},
	*/
}
