package interfaces

// ISink is an interface for sinks (https://rdf.js.org/stream-spec/#sink-interface).
type ISink interface {
	Import(IStream)
}
