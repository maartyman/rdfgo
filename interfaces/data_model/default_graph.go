package interfaces

// DefaultGraphType is an integer representing the default graph type.
const DefaultGraphType TermType = 4
const defaultGraphTypeString = "DefaultGraph"

// IDefaultGraph is an interface for default graphs (https://rdf.js.org/data-model-spec/#defaultgraph-interface).
type IDefaultGraph interface {
	ITerm
}
