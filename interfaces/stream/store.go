package interfaces

import (
	. "github.com/maartyman/rdfgo/interfaces/data_model"
)

// IStore is an interface for stores (https://rdf.js.org/stream-spec/#store-interface).
type IStore interface {
	ISource
	ISink
	Remove(IStream)
	RemoveMatches(ITerm, ITerm, ITerm, ITerm)
	DeleteGraph(ITerm)
}
