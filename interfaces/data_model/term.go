package interfaces

// ITerm is an interface for terms (https://rdf.js.org/data-model-spec/#term-interface).
type ITerm interface {
	Equals(other ITerm) bool
	GetType() TermType
	GetValue() string
	ToString() string
}

// TermType is an alias for an integer representing different terms.
type TermType int

var values = [...]string{
	namedNodeTypeString,
	literalTypeString,
	blankNodeTypeString,
	variableTypeString,
	defaultGraphTypeString,
	quadTypeString,
}

func (t TermType) String() string {
	return values[t]
}

func (t TermType) EnumIndex() int {
	return int(t)
}
