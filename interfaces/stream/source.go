package interfaces

import (
	. "rdfgo/interfaces/data_model"
)

type ISource interface {
	Match(ITerm, ITerm, ITerm, ITerm) IStream
}
