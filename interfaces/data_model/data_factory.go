package interfaces

// IDataFactory is an interface for data factories (https://rdf.js.org/data-model-spec/#datafactory-interface).
type IDataFactory interface {
	NamedNode(string) INamedNode
	BlankNode(string) IBlankNode
	Literal(string, string, INamedNode) ILiteral
	Variable(string) IVariable
	DefaultGraph() IDefaultGraph
	Quad(ITerm, ITerm, ITerm, ITerm) (IQuad, error)
}
