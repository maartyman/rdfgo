package interfaces

type ITerm interface {
	Equals(other ITerm) bool
	GetType() TermType
	GetValue() string
	ToString() string
}

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
