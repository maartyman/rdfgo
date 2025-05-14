package interfaces

import (
	. "github.com/maartyman/rdfgo/interfaces/data_model"
)

// IDatasetCoreFactory is an interface for dataset core factories (https://rdf.js.org/dataset-spec/#datasetcorefactory-interface).
type IDatasetCoreFactory interface {
	DatasetFromArray([]IQuad) IDataset
	DatasetFromDataset(dataset IDataset) IDataset
}
