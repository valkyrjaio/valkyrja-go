/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package header

import (
	"maps"
	"slices"
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

type HeaderCollection struct {
	headers map[string]contract.HeaderContract
}

// NewHeaderCollection builds a collection that holds each header. A second
// header of the same name replaces the first one.
func NewHeaderCollection(headers ...contract.HeaderContract) *HeaderCollection {
	collection := &HeaderCollection{headers: map[string]contract.HeaderContract{}}

	for _, header := range headers {
		collection.headers[header.GetNormalizedName()] = header
	}

	return collection
}

// Has reports whether the collection holds a header under the name.
func (c *HeaderCollection) Has(name string) bool {
	_, found := c.headers[strings.ToLower(name)]

	return found
}

// Get returns the header under the name. It reports a failure where the
// collection holds no header under the name.
func (c *HeaderCollection) Get(name string) (contract.HeaderContract, error) {
	header, found := c.headers[strings.ToLower(name)]
	if !found {
		return nil, exception.NewHttpHeaderInvalidHeaderNameError(name)
	}

	return header, nil
}

// GetHeaderLine returns every value of the header as one comma-separated string,
// and an empty string where the collection holds no header under the name.
func (c *HeaderCollection) GetHeaderLine(name string) string {
	header, err := c.Get(name)
	if err != nil {
		return ""
	}

	return header.GetHeaderLine()
}

// GetAll returns a copy of every header, keyed by its normalized name.
func (c *HeaderCollection) GetAll() map[string]contract.HeaderContract {
	return maps.Clone(c.headers)
}

// GetOnly returns the headers that the names identify.
func (c *HeaderCollection) GetOnly(names ...string) map[string]contract.HeaderContract {
	return c.filterByName(names, true)
}

// GetAllExcept returns every header that the names do not identify.
func (c *HeaderCollection) GetAllExcept(names ...string) map[string]contract.HeaderContract {
	return c.filterByName(names, false)
}

// WithHeader returns a copy of the collection with the header in it. A header of
// the same name is replaced rather than merged.
func (c *HeaderCollection) WithHeader(header contract.HeaderContract) contract.HeaderCollectionContract {
	copied := c.clone()
	copied.headers[header.GetNormalizedName()] = header

	return copied
}

// WithoutHeader returns a copy of the collection without the header.
func (c *HeaderCollection) WithoutHeader(name string) contract.HeaderCollectionContract {
	copied := c.clone()

	delete(copied.headers, strings.ToLower(name))

	return copied
}

// WithHeaders returns a copy of the collection that holds the headers and
// nothing else.
func (c *HeaderCollection) WithHeaders(headers ...contract.HeaderContract) contract.HeaderCollectionContract {
	return NewHeaderCollection(headers...)
}

// WithAddedHeaders returns a copy of the collection with the headers added to
// the ones it holds.
func (c *HeaderCollection) WithAddedHeaders(headers ...contract.HeaderContract) contract.HeaderCollectionContract {
	copied := c.clone()

	for _, header := range headers {
		name := header.GetNormalizedName()

		existing, found := copied.headers[name]
		if !found {
			copied.headers[name] = header

			continue
		}

		copied.headers[name] = existing.WithAddedValues(header.GetValues()...)
	}

	return copied
}

// filterByName returns the headers whose normalized name is in the names, where
// keep is true, and the headers whose name is not, where keep is false.
func (c *HeaderCollection) filterByName(names []string, keep bool) map[string]contract.HeaderContract {
	normalized := make([]string, 0, len(names))

	for _, name := range names {
		normalized = append(normalized, strings.ToLower(name))
	}

	filtered := map[string]contract.HeaderContract{}

	for name, header := range c.headers {
		if slices.Contains(normalized, name) == keep {
			filtered[name] = header
		}
	}

	return filtered
}

// clone returns a copy of the collection that shares no map with it.
func (c *HeaderCollection) clone() *HeaderCollection {
	return &HeaderCollection{headers: maps.Clone(c.headers)}
}
