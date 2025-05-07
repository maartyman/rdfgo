package rdfgo

import (
	"github.com/maartyman/rdfgo/interfaces"
	datamodel "github.com/maartyman/rdfgo/lib/data_model"
	parser "github.com/maartyman/rdfgo/lib/parser"
	stream "github.com/maartyman/rdfgo/lib/stream"
)

// Interfaces exports
// Data Model
type (
	IBlankNode    = interfaces.IBlankNode
	IDataFactory  = interfaces.IDataFactory
	IDefaultGraph = interfaces.IDefaultGraph
	ILiteral      = interfaces.ILiteral
	INamedNode    = interfaces.INamedNode
	IQuad         = interfaces.IQuad
	ITerm         = interfaces.ITerm
	TermType      = interfaces.TermType
	IVariable     = interfaces.IVariable
)

// Dataset
type (
	IDataset            = interfaces.IDataset
	IDatasetFactory     = interfaces.IDatasetFactory
	IDatasetCore        = interfaces.IDatasetCore
	IDatasetCoreFactory = interfaces.IDatasetCoreFactory
)

// Stream
type (
	ISink   = interfaces.ISink
	ISource = interfaces.ISource
	IStore  = interfaces.IStore
	IStream = interfaces.IStream
)

// Term types
const (
	BlankNodeType    = interfaces.BlankNodeType
	DefaultGraphType = interfaces.DefaultGraphType
	LiteralType      = interfaces.LiteralType
	NamedNodeType    = interfaces.NamedNodeType
	QuadType         = interfaces.QuadType
	VariableType     = interfaces.VariableType
)

// Data Model exports
// Term constructors
var (
	NewBlankNode      = datamodel.NewBlankNode
	NewBooleanLiteral = datamodel.NewBooleanLiteral
	NewDataFactory    = datamodel.NewDataFactory
	NewDecimalLiteral = datamodel.NewDecimalLiteral
	NewDefaultGraph   = datamodel.NewDefaultGraph
	NewDoubleLiteral  = datamodel.NewDoubleLiteral
	NewIntegerLiteral = datamodel.NewIntegerLiteral
	NewLiteral        = datamodel.NewLiteral
	NewNamedNode      = datamodel.NewNamedNode
	NewQuad           = datamodel.NewQuad
	NewStringLiteral  = datamodel.NewStringLiteral
	NewVariable       = datamodel.NewVariable
)

// Other
var (
	IRI                    = datamodel.IRI
	SubjectTermTypeError   = datamodel.SubjectTermTypeError
	PredicateTermTypeError = datamodel.PredicateTermTypeError
	ObjectTermTypeError    = datamodel.ObjectTermTypeError
	GraphTermTypeError     = datamodel.GraphTermTypeError
	DefaultGraphValue      = datamodel.DefaultGraphValue
	DefaultGraphString     = datamodel.DefaultGraphString
	QuadValue              = datamodel.QuadValue
)

type DataFactory = datamodel.DataFactory

// Parser exports
var (
	Parse     = parser.Parse
	ParseFile = parser.ParseFile
)

type Options = parser.Options

// Stream exports
var (
	NewStream     = stream.NewStream
	ArrayToStream = stream.ArrayToStream
	NewStore      = stream.NewStore
)

type Store = stream.Store
type Stream = stream.Stream
