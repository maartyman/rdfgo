package interfaces

// BlankNodeType is an integer representing the blank node type.
const BlankNodeType TermType = 2
const blankNodeTypeString = "BlankNode"

// IBlankNode is an interface for blank nodes (https://rdf.js.org/data-model-spec/#blanknode-interface).
type IBlankNode interface {
	ITerm
}
