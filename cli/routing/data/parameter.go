/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package data holds the commands of the CLI router, and the parameters that a
// command takes.
package data

import (
	"slices"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
)

type parameter struct {
	name        string
	description string
	cast        contract.CastFunc
}

// GetName returns the name of the parameter.
func (p *parameter) GetName() string {
	return p.name
}

// HasCast reports whether the parameter casts its values to a type.
func (p *parameter) HasCast() bool {
	return p.cast != nil
}

// GetCast returns what converts the values of the parameter, and nil where the
// parameter casts nothing.
func (p *parameter) GetCast() contract.CastFunc {
	return p.cast
}

// GetDescription returns the description that the help text prints.
func (p *parameter) GetDescription() string {
	return p.description
}

// getCastValues returns each value, cast to the type that the parameter names.
func getCastValues(cast contract.CastFunc, values []string) ([]any, error) {
	castValues := make([]any, 0, len(values))

	for _, value := range values {
		if cast == nil {
			castValues = append(castValues, value)

			continue
		}

		castValue, err := cast(value)
		if err != nil {
			return nil, err
		}

		castValues = append(castValues, castValue)
	}

	return castValues, nil
}

// appendUnique returns the values with each added value that the values do not
// hold already.
func appendUnique(values []string, added []string) []string {
	combined := make([]string, 0, len(values)+len(added))
	combined = append(combined, values...)

	for _, value := range added {
		if slices.Contains(combined, value) {
			continue
		}

		combined = append(combined, value)
	}

	return combined
}
