package interfaces

import (
	. "github.com/maartyman/rdfgo/interfaces/data_model"
)

// IStream is an interface for streams (https://rdf.js.org/stream-spec/#stream-interface).
type IStream chan IQuad
