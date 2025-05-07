package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
)

type dataFactory struct {
	blankNodeCounter int
}

type DataFactory interface {
	interfaces.IDataFactory
	SimpleLiteral(value string) interfaces.ILiteral
}

func NewDataFactory() DataFactory {
	return &dataFactory{0}
}

func (df *dataFactory) NamedNode(value string) interfaces.INamedNode {
	return NewNamedNode(value)
}

func (df *dataFactory) BlankNode(value string) interfaces.IBlankNode {
	return NewBlankNode(value)
}

func (df *dataFactory) SimpleLiteral(value string) interfaces.ILiteral {
	return NewLiteral(value, "", IRI.XSD.String)
}

func (df *dataFactory) Literal(value string, language string, datatype interfaces.INamedNode) interfaces.ILiteral {
	if datatype == nil {
		datatype = IRI.XSD.String
	}
	return NewLiteral(value, language, datatype)
}

func (df *dataFactory) Variable(value string) interfaces.IVariable {
	return NewVariable(value)
}

func (df *dataFactory) DefaultGraph() interfaces.IDefaultGraph {
	return NewDefaultGraph()
}

func (df *dataFactory) Quad(
	subject interfaces.ITerm,
	predicate interfaces.ITerm,
	object interfaces.ITerm,
	graph interfaces.ITerm,
) (interfaces.IQuad, error) {
	return NewQuad(subject, predicate, object, graph)
}
