package interfaces

import (
	. "rdfgo/interfaces/data_model"
)

type IStore interface {
	ISource
	ISink
	Remove(IStream)
	RemoveMatches(ITerm, ITerm, ITerm, ITerm)
	DeleteGraph(ITerm)
}
