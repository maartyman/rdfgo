package interfaces

// NamedNodeType is an integer representing the named node type.
const NamedNodeType TermType = 0
const namedNodeTypeString = "NamedNode"

// INamedNode is an interface for named nodes (https://rdf.js.org/data-model-spec/#namednode-interface).
type INamedNode interface {
	ITerm
}
