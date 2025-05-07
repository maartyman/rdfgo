// Package rdfgo provides a Go implementation of the RDF.js data model, dataset, and stream interfaces together with their implementations.
//
// This package consists of interfaces in the `interfaces` subpackage and implementations in the `lib` subpackage.
//
// ./interfaces:
//
// - data model: blank node, literal, named node, quad, term, variable, default graph, and data factory.
//
// - dataset: datasets, dataset core, dataset core factory, and dataset factory.
//
// - stream: sink, source, store, and stream.
//
// ./lib (implementations):
//
// - data model: blank node, literal, named node, quad, term, variable, default graph, and data factory.
//
// - parser: parser for n-quads (and n-turtle) and turtle formats.
//
// - stream: store, and stream.
//
// Examples can be found on the GitHub page: https://github.com/maartyman/rdfgo
package rdfgo
