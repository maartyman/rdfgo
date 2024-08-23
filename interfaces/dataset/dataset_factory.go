package interfaces

import (
	. "github.com/maartyman/rdfgo/interfaces/data_model"
)

type IDatasetFactory interface {
	DatasetFromArray([]IQuad) IDataset
	DatasetFromDataset(dataset IDataset) IDataset
}
