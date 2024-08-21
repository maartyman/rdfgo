package interfaces

const QuadType TermType = 5
const quadTypeString = "Quad"

type IQuad interface {
	ITerm
	GetSubject() ITerm
	GetPredicate() ITerm
	GetObject() ITerm
	GetGraph() ITerm
}
