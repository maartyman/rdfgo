package interfaces

import datamodel "github.com/maartyman/rdfgo/interfaces/data_model"
import dataset "github.com/maartyman/rdfgo/interfaces/dataset"
import stream "github.com/maartyman/rdfgo/interfaces/stream"

type IBlankNode = datamodel.IBlankNode
type IDataFactory = datamodel.IDataFactory
type IDefaultGraph = datamodel.IDefaultGraph
type ILiteral = datamodel.ILiteral
type INamedNode = datamodel.INamedNode
type IQuad = datamodel.IQuad
type ITerm = datamodel.ITerm
type TermType = datamodel.TermType
type IVariable = datamodel.IVariable

const BlankNodeType = datamodel.BlankNodeType
const DefaultGraphType = datamodel.DefaultGraphType
const LiteralType = datamodel.LiteralType
const NamedNodeType = datamodel.NamedNodeType
const QuadType = datamodel.QuadType
const VariableType = datamodel.VariableType

type IDataset = dataset.IDataset
type IDatasetFactory = dataset.IDatasetFactory
type IDatasetCore = dataset.IDatasetCore
type IDatasetCoreFactory = dataset.IDatasetCoreFactory

type ISink = stream.ISink
type ISource = stream.ISource
type IStore = stream.IStore
type IStream = stream.IStream
