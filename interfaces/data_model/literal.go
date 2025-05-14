package interfaces

// LiteralType is an integer representing the literal type.
const LiteralType TermType = 1
const literalTypeString = "Literal"

// ILiteral is an interface for literals (https://rdf.js.org/data-model-spec/#literal-interface).
type ILiteral interface {
	ITerm
	GetLanguage() string
	GetDatatype() INamedNode
}
