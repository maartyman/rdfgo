package interfaces

const LiteralType TermType = 1
const literalTypeString = "Literal"

type ILiteral interface {
	ITerm
	GetLanguage() string
	GetDatatype() INamedNode
}
