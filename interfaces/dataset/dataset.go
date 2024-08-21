package interfaces

import (
	. "rdfgo/interfaces/data_model"
	. "rdfgo/interfaces/stream"
)

type IDataset interface {
	IDatasetCore
	AddAll(IDataset) IDatasetCore
	Contains(IDataset) bool
	DeleteMatches(ITerm, ITerm, ITerm, ITerm) IDatasetCore
	Difference(IDataset) IDatasetCore
	Equals(IDataset) bool
	Every(func(IQuad) bool) bool
	Filter(func(IQuad) bool) IDatasetCore
	ForEach(func(IQuad))
	Import(IStream) IDatasetCore
	Intersection(IDataset) IDatasetCore
	MapQuads(func(IQuad) IQuad) IDatasetCore
	Reduce(func(interface{}, IQuad) interface{}, interface{}) interface{}
	Some(func(IQuad) bool) bool
	ToArray() []IQuad
	ToCanonical() string
	ToStream() IStream
	ToString() string
	Union(IDataset) IDatasetCore
}
