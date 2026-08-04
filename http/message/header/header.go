/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package header

import (
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
)

// nameDeliminator separates the name of a header from its values.
const nameDeliminator = ":"

// headerValueDeliminator separates one value of a header from the next.
const headerValueDeliminator = ","

// Header is one header of a message.
type Header struct {
	name           string
	normalizedName string
	values         []contract.ValueContract
}

// NewHeader builds a header from a name and its values. It reports a failure
// where the name is not one that a header carries.
func NewHeader(name string, values ...contract.ValueContract) (*Header, error) {
	err := ValidateName(name)
	if err != nil {
		return nil, err
	}

	return &Header{
		name:           name,
		normalizedName: strings.ToLower(name),
		values:         values,
	}, nil
}

// NewHeaderFromValue builds a header by reading the name and each value out of
// one header line.
func NewHeaderFromValue(line string) (*Header, error) {
	name, rendered, found := strings.Cut(line, nameDeliminator)
	if !found {
		name = line
		rendered = ""
	}

	parts := strings.Split(rendered, headerValueDeliminator)

	values := make([]contract.ValueContract, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)

		err := ValidateValue(trimmed)
		if err != nil {
			return nil, err
		}

		values = append(values, value.NewValueFromValue(trimmed))
	}

	return NewHeader(strings.TrimSpace(name), values...)
}

// GetName returns the name of the header, as the sender wrote it.
func (h *Header) GetName() string {
	return h.name
}

// GetNormalizedName returns the name of the header in lower case, which is the
// form that the collection keys on.
func (h *Header) GetNormalizedName() string {
	return h.normalizedName
}

// WithName returns a copy of the header under another name.
//
// The other ports throw where the name is invalid. A `With` method returns the
// contract in every port, so this one keeps the receiver's name instead, and a
// caller that needs the failure builds the header with `NewHeader`.
func (h *Header) WithName(name string) contract.HeaderContract {
	if !IsValidName(name) {
		return h
	}

	copied := *h
	copied.name = name
	copied.normalizedName = strings.ToLower(name)

	return &copied
}

// GetValues returns each value of the header.
func (h *Header) GetValues() []contract.ValueContract {
	return h.values
}

// WithValues returns a copy of the header with other values.
func (h *Header) WithValues(values ...contract.ValueContract) contract.HeaderContract {
	copied := *h
	copied.values = values

	return &copied
}

// WithAddedValues returns a copy of the header with the values added after the
// ones it holds.
func (h *Header) WithAddedValues(values ...contract.ValueContract) contract.HeaderContract {
	copied := *h

	combined := make([]contract.ValueContract, 0, len(h.values)+len(values))
	combined = append(combined, h.values...)
	combined = append(combined, values...)

	copied.values = combined

	return &copied
}

// GetHeaderLine returns every value of the header as one comma-separated string.
func (h *Header) GetHeaderLine() string {
	parts := make([]string, 0, len(h.values))

	for _, headerValue := range h.values {
		rendered := strings.TrimSpace(headerValue.String())
		if rendered == "" {
			continue
		}

		parts = append(parts, rendered)
	}

	return strings.Join(parts, headerValueDeliminator+" ")
}

// String returns the whole header as a string, and an empty string where the
// header carries no value.
func (h *Header) String() string {
	line := h.GetHeaderLine()
	if line == "" {
		return ""
	}

	return h.name + nameDeliminator + " " + line
}
