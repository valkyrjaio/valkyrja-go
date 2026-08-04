/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package structure

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

// ResponseStructure shapes the data that a response returns.
//
// The structure maps the name a field carries inside the application to the name
// that a client reads, so a rename inside never reaches a client.
type ResponseStructure struct {
	name   string
	value  any
	fields map[string]string
}

// NewResponseStructure builds a structure that renames each field, keyed by the
// name inside the application and valued by the name that a client reads.
func NewResponseStructure(name string, fields map[string]string) *ResponseStructure {
	return &ResponseStructure{
		name:   name,
		fields: fields,
	}
}

// GetName returns the name of the structure.
func (s *ResponseStructure) GetName() string {
	return s.name
}

// GetValue returns the value of the structure.
func (s *ResponseStructure) GetValue() any {
	return s.value
}

// WithValue returns a copy of the structure with another value.
func (s *ResponseStructure) WithValue(value any) *ResponseStructure {
	copied := *s
	copied.value = value

	return &copied
}

// GetFields returns each field that the structure renames.
func (s *ResponseStructure) GetFields() map[string]string {
	return s.fields
}

// GetStructuredData returns the data in the shape that the structure names.
//
// A field that the data does not carry reads as nil where includeAll is true,
// and it is left out where includeAll is false. A field that the structure does
// not name never reaches a client, whether the data carries it or not.
func (s *ResponseStructure) GetStructuredData(data map[string]any, includeAll bool) map[string]any {
	structured := make(map[string]any, len(s.fields))

	for field, renamed := range s.fields {
		held, carried := data[field]

		if !includeAll && !carried {
			continue
		}

		structured[renamed] = held
	}

	return structured
}

// A response structure satisfies its contract, which the compiler checks at
// build time rather than at run time.
var _ contract.ResponseStructContract = (*ResponseStructure)(nil)
