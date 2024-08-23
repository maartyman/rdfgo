package rdfgo

import (
	"rdfgo/interfaces"
)

type DataFactory struct {
	blankNodeCounter int
}

func NewDataFactory() *DataFactory {
	return &DataFactory{0}
}

func (df *DataFactory) NamedNode(value string) interfaces.INamedNode {
	return NewNamedNode(value)
}

func (df *DataFactory) BlankNode(value string) interfaces.IBlankNode {
	return NewBlankNode(value)
}

func (df *DataFactory) SimpleLiteral(value string) interfaces.ILiteral {
	return NewLiteral(value, "", df.NamedNode("http://www.w3.org/2001/XMLSchema#string"))
}

func (df *DataFactory) Literal(value string, language string, datatype interfaces.INamedNode) interfaces.ILiteral {
	if datatype == nil {
		datatype = df.NamedNode("http://www.w3.org/2001/XMLSchema#string")
	}
	return NewLiteral(value, language, datatype)
}

func (df *DataFactory) Variable(value string) interfaces.IVariable {
	return NewVariable(value)
}

func (df *DataFactory) DefaultGraph() interfaces.IDefaultGraph {
	return NewDefaultGraph()
}

func (df *DataFactory) Quad(
	subject interfaces.ITerm,
	predicate interfaces.ITerm,
	object interfaces.ITerm,
	graph interfaces.ITerm,
) (interfaces.IQuad, error) {
	return NewQuad(subject, predicate, object, graph)
}
