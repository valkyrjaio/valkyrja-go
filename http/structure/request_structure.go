/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package structure names the fields that a request carries and that a response
// returns.
//
// The other ports spell this segment `Struct`. `struct` is a Go keyword, so it
// cannot be a package name, and this port spells the word in full.
package structure

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

// SourceFunc returns the parameters of a request that a struct reads.
//
// The other ports declare an abstract `RequestStruct` and override the two
// methods that read a request. Go has no abstract type and no method override,
// so a struct holds the function that names the collection to read.
type SourceFunc func(request contract.ServerRequestContract) contract.ParamCollectionContract

// RequestStructure reads the fields that it names out of a request.
type RequestStructure struct {
	name   string
	value  any
	fields []string
	source SourceFunc
}

// NewRequestStructure builds a structure that reads the fields from the
// collection that the source names.
func NewRequestStructure(name string, fields []string, source SourceFunc) *RequestStructure {
	return &RequestStructure{
		name:   name,
		fields: fields,
		source: source,
	}
}

// NewQueryRequestStructure builds a structure that reads the query parameters.
func NewQueryRequestStructure(name string, fields ...string) *RequestStructure {
	return NewRequestStructure(name, fields, func(
		request contract.ServerRequestContract,
	) contract.ParamCollectionContract {
		return request.GetQueryParams()
	})
}

// NewParsedBodyRequestStructure builds a structure that reads the parsed body.
func NewParsedBodyRequestStructure(name string, fields ...string) *RequestStructure {
	return NewRequestStructure(name, fields, func(
		request contract.ServerRequestContract,
	) contract.ParamCollectionContract {
		return request.GetParsedBody()
	})
}

// NewJsonRequestStructure builds a structure that reads the parsed JSON body.
//
// Warning: a request that parsed no JSON body reads as empty rather than as a
// failure. The other ports throw where the request is not a JSON request; every
// server request of this port carries the parsed JSON, so there is no other type
// to be, and an empty collection is what a request with no JSON body holds.
func NewJsonRequestStructure(name string, fields ...string) *RequestStructure {
	return NewRequestStructure(name, fields, func(
		request contract.ServerRequestContract,
	) contract.ParamCollectionContract {
		return request.GetParsedJson()
	})
}

// GetName returns the name of the structure.
func (s *RequestStructure) GetName() string {
	return s.name
}

// GetValue returns the value of the structure.
func (s *RequestStructure) GetValue() any {
	return s.value
}

// WithValue returns a copy of the structure with another value.
func (s *RequestStructure) WithValue(value any) *RequestStructure {
	copied := *s
	copied.value = value

	return &copied
}

// GetFields returns each field that the structure names.
func (s *RequestStructure) GetFields() []string {
	return s.fields
}

// GetDataFromRequest returns the fields that the structure names, read from the
// request.
func (s *RequestStructure) GetDataFromRequest(request contract.ServerRequestContract) map[string]any {
	return s.source(request).GetOnly(s.fields...)
}

// DetermineIfRequestContainsExtraData reports whether the request carries a
// field that the structure does not name.
func (s *RequestStructure) DetermineIfRequestContainsExtraData(request contract.ServerRequestContract) bool {
	return len(s.source(request).GetAllExcept(s.fields...)) > 0
}

// A request structure satisfies its contract, which the compiler checks at build
// time rather than at run time.
var _ contract.RequestStructContract = (*RequestStructure)(nil)
