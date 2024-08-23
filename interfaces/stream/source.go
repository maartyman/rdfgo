package interfaces

import (
	. "github.com/maartyman/rdfgo/interfaces/data_model"
)

type ISource interface {
	Match(ITerm, ITerm, ITerm, ITerm) IStream
}
