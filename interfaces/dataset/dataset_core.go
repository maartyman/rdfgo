package interfaces

import (
	. "github.com/maartyman/rdfgo/interfaces/data_model"
)

// IDatasetCore is an interface for dataset cores (https://rdf.js.org/dataset-spec/#datasetcore-interface).
type IDatasetCore interface {
	GetSize() int
	Add(IQuad) IDatasetCore
	Delete(IQuad) IDatasetCore
	Has(IQuad) bool
	Match(ITerm, ITerm, ITerm, ITerm) IDatasetCore
}
