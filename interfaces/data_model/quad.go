package interfaces

// QuadType is an integer representing the quad type.
const QuadType TermType = 5
const quadTypeString = "Quad"

// IQuad is an interface for quads (https://rdf.js.org/data-model-spec/#quad-interface).
type IQuad interface {
	ITerm
	GetSubject() ITerm
	GetPredicate() ITerm
	GetObject() ITerm
	GetGraph() ITerm
}
