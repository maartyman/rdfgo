package interfaces

import (
	. "github.com/maartyman/rdfgo/interfaces/data_model"
)

// ISource is an interface for sources (https://rdf.js.org/stream-spec/#source-interface).
type ISource interface {
	Match(ITerm, ITerm, ITerm, ITerm) IStream
}
