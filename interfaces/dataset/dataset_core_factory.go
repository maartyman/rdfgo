package interfaces

import (
	. "github.com/maartyman/rdfgo/interfaces/data_model"
)

type IDatasetCoreFactory interface {
	DatasetFromArray([]IQuad) IDataset
	DatasetFromDataset(dataset IDataset) IDataset
}
