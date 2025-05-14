package interfaces

import datamodel "github.com/maartyman/rdfgo/interfaces/data_model"
import dataset "github.com/maartyman/rdfgo/interfaces/dataset"
import stream "github.com/maartyman/rdfgo/interfaces/stream"

// IBlankNode is an interface for blank nodes (https://rdf.js.org/data-model-spec/#blanknode-interface).
type IBlankNode = datamodel.IBlankNode

// IDataFactory is an interface for data factories (https://rdf.js.org/data-model-spec/#datafactory-interface).
type IDataFactory = datamodel.IDataFactory

// IDefaultGraph is an interface for default graphs (https://rdf.js.org/data-model-spec/#defaultgraph-interface).
type IDefaultGraph = datamodel.IDefaultGraph

// ILiteral is an interface for literals (https://rdf.js.org/data-model-spec/#literal-interface).
type ILiteral = datamodel.ILiteral

// INamedNode is an interface for named nodes (https://rdf.js.org/data-model-spec/#namednode-interface).
type INamedNode = datamodel.INamedNode

// IQuad is an interface for quads (https://rdf.js.org/data-model-spec/#quad-interface).
type IQuad = datamodel.IQuad

// ITerm is an interface for terms (https://rdf.js.org/data-model-spec/#term-interface).
type ITerm = datamodel.ITerm

// TermType is an alias for an integer representing different terms.
type TermType = datamodel.TermType

// IVariable is an interface for variables (https://rdf.js.org/data-model-spec/#variable-interface).
type IVariable = datamodel.IVariable

// BlankNodeType is an integer representing the blank node type.
const BlankNodeType = datamodel.BlankNodeType

// DefaultGraphType is an integer representing the default graph type.
const DefaultGraphType = datamodel.DefaultGraphType

// LiteralType is an integer representing the literal type.
const LiteralType = datamodel.LiteralType

// NamedNodeType is an integer representing the named node type.
const NamedNodeType = datamodel.NamedNodeType

// QuadType is an integer representing the quad type.
const QuadType = datamodel.QuadType

// VariableType is an integer representing the variable type.
const VariableType = datamodel.VariableType

// IDataset is an interface for datasets (https://rdf.js.org/dataset-spec/#dataset-interface).
type IDataset = dataset.IDataset

// IDatasetFactory is an interface for dataset factories (https://rdf.js.org/dataset-spec/#datasetfactory-interface).
type IDatasetFactory = dataset.IDatasetFactory

// IDatasetCore is an interface for dataset cores (https://rdf.js.org/dataset-spec/#datasetcore-interface).
type IDatasetCore = dataset.IDatasetCore

// IDatasetCoreFactory is an interface for dataset core factories (https://rdf.js.org/dataset-spec/#datasetcorefactory-interface).
type IDatasetCoreFactory = dataset.IDatasetCoreFactory

// ISink is an interface for sinks (https://rdf.js.org/stream-spec/#sink-interface).
type ISink = stream.ISink

// ISource is an interface for sources (https://rdf.js.org/stream-spec/#source-interface).
type ISource = stream.ISource

// IStore is an interface for stores (https://rdf.js.org/stream-spec/#store-interface).
type IStore = stream.IStore

// IStream is an interface for streams (https://rdf.js.org/stream-spec/#stream-interface).
type IStream = stream.IStream
