package interfaces

type IDataFactory interface {
	NamedNode(string) INamedNode
	BlankNode(string) IBlankNode
	Literal(string, string, INamedNode) ILiteral
	Variable(string) IVariable
	DefaultGraph() IDefaultGraph
	Quad(ITerm, ITerm, ITerm, ITerm) IQuad
}
