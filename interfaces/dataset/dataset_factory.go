package interfaces

import (
	. "github.com/maartyman/rdfgo/interfaces/data_model"
)

// IDatasetFactory is an interface for dataset factories (https://rdf.js.org/dataset-spec/#datasetfactory-interface).
type IDatasetFactory interface {
	DatasetFromArray([]IQuad) IDataset
	DatasetFromDataset(dataset IDataset) IDataset
}
