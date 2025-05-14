package interfaces

// VariableType is an integer representing the variable type.
const VariableType TermType = 3
const variableTypeString = "Variable"

// IVariable is an interface for variables (https://rdf.js.org/data-model-spec/#variable-interface).
type IVariable interface {
	ITerm
}
