package interfaces

import (
	. "rdfgo/interfaces/data_model"
)

type IDatasetCore interface {
	GetSize() int
	Add(IQuad) IDatasetCore
	Delete(IQuad) IDatasetCore
	Has(IQuad) bool
	Match(ITerm, ITerm, ITerm, ITerm) IDatasetCore
}
