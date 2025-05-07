package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
	datamodel "github.com/maartyman/rdfgo/lib/data_model"
	parser "github.com/maartyman/rdfgo/lib/parser"
	stream "github.com/maartyman/rdfgo/lib/stream"
)

// Interfaces types:

// IBlankNode is an interface for blank nodes (https://rdf.js.org/data-model-spec/#blanknode-interface).
type IBlankNode = interfaces.IBlankNode

// IDataFactory is an interface for data factories (https://rdf.js.org/data-model-spec/#datafactory-interface).
type IDataFactory = interfaces.IDataFactory

// IDefaultGraph is an interface for default graphs (https://rdf.js.org/data-model-spec/#defaultgraph-interface).
type IDefaultGraph = interfaces.IDefaultGraph

// ILiteral is an interface for literals (https://rdf.js.org/data-model-spec/#literal-interface).
type ILiteral = interfaces.ILiteral

// INamedNode is an interface for named nodes (https://rdf.js.org/data-model-spec/#namednode-interface).
type INamedNode = interfaces.INamedNode

// IQuad is an interface for quads (https://rdf.js.org/data-model-spec/#quad-interface).
type IQuad = interfaces.IQuad

// ITerm is an interface for terms (https://rdf.js.org/data-model-spec/#term-interface).
type ITerm = interfaces.ITerm

// TermType is an alias for an integer representing different terms.
type TermType = interfaces.TermType

// IVariable is an interface for variables (https://rdf.js.org/data-model-spec/#variable-interface).
type IVariable = interfaces.IVariable

// IDataset is an interface for datasets (https://rdf.js.org/dataset-spec/#dataset-interface).
type IDataset = interfaces.IDataset

// IDatasetFactory is an interface for dataset factories (https://rdf.js.org/dataset-spec/#datasetfactory-interface).
type IDatasetFactory = interfaces.IDatasetFactory

// IDatasetCore is an interface for dataset cores (https://rdf.js.org/dataset-spec/#datasetcore-interface).
type IDatasetCore = interfaces.IDatasetCore

// IDatasetCoreFactory is an interface for dataset core factories (https://rdf.js.org/dataset-spec/#datasetcorefactory-interface).
type IDatasetCoreFactory = interfaces.IDatasetCoreFactory

// ISink is an interface for sinks (https://rdf.js.org/stream-spec/#sink-interface).
type ISink = interfaces.ISink

// ISource is an interface for sources (https://rdf.js.org/stream-spec/#source-interface).
type ISource = interfaces.ISource

// IStore is an interface for stores (https://rdf.js.org/stream-spec/#store-interface).
type IStore = interfaces.IStore

// IStream is an interface for streams (https://rdf.js.org/stream-spec/#stream-interface).
type IStream = interfaces.IStream

// Interfaces constants:

// BlankNodeType is an integer representing the blank node type.
const BlankNodeType = interfaces.BlankNodeType

// DefaultGraphType is an integer representing the default graph type.
const DefaultGraphType = interfaces.DefaultGraphType

// LiteralType is an integer representing the literal type.
const LiteralType = interfaces.LiteralType

// NamedNodeType is an integer representing the named node type.
const NamedNodeType = interfaces.NamedNodeType

// QuadType is an integer representing the quad type.
const QuadType = interfaces.QuadType

// VariableType is an integer representing the variable type.
const VariableType = interfaces.VariableType

// Data model functions:

// NewBlankNode is a constructor for creating a blank node. This returns an implementation of the interfaces.IBlankNode interface.
var NewBlankNode = datamodel.NewBlankNode

// NewBooleanLiteral is a constructor for creating a boolean literal. This returns an implementation of the interfaces.ILiteral interface.
var NewBooleanLiteral = datamodel.NewBooleanLiteral

// NewDataFactory is a constructor for creating a data factory. This returns an implementation of the DataFactory interface.
var NewDataFactory = datamodel.NewDataFactory

// NewDecimalLiteral is a constructor for creating a decimal literal. This returns an implementation of the interfaces.ILiteral interface.
var NewDecimalLiteral = datamodel.NewDecimalLiteral

// NewDefaultGraph is a constructor for creating a default graph. This returns an implementation of the interfaces.IDefaultGraph interface.
var NewDefaultGraph = datamodel.NewDefaultGraph

// NewDoubleLiteral is a constructor for creating a double literal. This returns an implementation of the interfaces.ILiteral interface.
var NewDoubleLiteral = datamodel.NewDoubleLiteral

// NewIntegerLiteral is a constructor for creating an integer literal. This returns an implementation of the interfaces.ILiteral interface.
var NewIntegerLiteral = datamodel.NewIntegerLiteral

// NewLiteral is a constructor for creating a literal. This returns an implementation of the interfaces.ILiteral interface.
var NewLiteral = datamodel.NewLiteral

// NewNamedNode is a constructor for creating a named node. This returns an implementation of the interfaces.INamedNode interface.
var NewNamedNode = datamodel.NewNamedNode

// NewQuad is a constructor for creating a quad. This returns an implementation of the interfaces.IQuad interface.
var NewQuad = datamodel.NewQuad

// NewStringLiteral is a constructor for creating a string literal. This returns an implementation of the interfaces.ILiteral interface.
var NewStringLiteral = datamodel.NewStringLiteral

// NewVariable is a constructor for creating a variable. This returns an implementation of the interfaces.IVariable interface.
var NewVariable = datamodel.NewVariable

// Data model variables:

// IRI gives access to common the IRI namespaces, like XSD, and RDF.
var IRI = datamodel.IRI

// SubjectTermTypeError is an error type that is thrown when a quad's subject is not of the correct type.
var SubjectTermTypeError = datamodel.SubjectTermTypeError

// PredicateTermTypeError is an error type that is thrown when a quad's predicate is not of the correct type.
var PredicateTermTypeError = datamodel.PredicateTermTypeError

// ObjectTermTypeError is an error type that is thrown when a quad's object is not of the correct type.
var ObjectTermTypeError = datamodel.ObjectTermTypeError

// GraphTermTypeError is an error type that is thrown when a quad's graph is not of the correct type.
var GraphTermTypeError = datamodel.GraphTermTypeError

// DefaultGraphValue is a constant that represents the default graph value.
var DefaultGraphValue = datamodel.DefaultGraphValue

// DefaultGraphString is a constant that represents the default graph string.
var DefaultGraphString = datamodel.DefaultGraphString

// QuadValue is a constant that represents the quad value. This is never used as the value of a quad is the combination of the subject, predicate, object, and graph.
var QuadValue = datamodel.QuadValue

// Data model types:

// DataFactory is an extension of the interfaces.IDataFactory interface. It includes a SimpleLiteral method that creates a literal with a default datatype of xsd:string.
type DataFactory = datamodel.DataFactory

// Parser functions:

// Parse is a function that parses a string into a stream of quads. It accepts an io.Reader and parser.Options as parameters, and returns a channel of quads and a channel of errors. If no format is specified, it defaults to turtle format.
var Parse = parser.Parse

// ParseFile is a function that parses a file into a stream of quads. It accepts a file path and parser.Options as parameters, and returns a channel of quads and a channel of errors. It supports n-triples, n-quads, and turtle formats (.nt, .nq, .ttl).
var ParseFile = parser.ParseFile

// Parser types

// Options is used to configure the parser. It includes options for the base IRI, and the format type.
type Options = parser.Options

// Stream functions:

// NewStream creates a new chanel of IQuad. It has some helper functions to combine streams, convert to an array, create a store, count the quads.
var NewStream = stream.NewStream

// ArrayToStream converts an array of quads to a Stream. It accepts a slice of interfaces.IQuad and returns a Stream.
var ArrayToStream = stream.ArrayToStream

// NewStore creates a new store of quads. A store indexes the quads by their subject, predicate, object, and graph. It accepts a slice of interfaces.IQuad and returns a Store.
var NewStore = stream.NewStore

// Stream types:

// Store is an extension of the interfaces.IStore interface. It has various methods to manipulate the store, like adding and removing quads, checking if a quad exists, and iterating over the quads.
type Store = stream.Store

// Stream is an extension of the interfaces.IStream interface. It has various methods to manipulate the stream, like counting the quads, importing a stream, and converting to an array.
type Stream = stream.Stream
