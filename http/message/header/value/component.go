/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package value holds one value of an HTTP header, and the parts that a value is
// built from.
package value

import (
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

// componentDeliminator separates the token of a component from its text.
const componentDeliminator = "="

// valueDeliminator separates one component of a value from the next.
const valueDeliminator = ";"

// Component is one token-and-text pair inside a header value.
type Component struct {
	token string
	text  string
}

// NewComponent builds a component from a token and a text. It trims each one,
// because a header carries a space after each deliminator.
func NewComponent(token string, text string) *Component {
	return &Component{
		token: strings.TrimSpace(token),
		text:  strings.TrimSpace(text),
	}
}

// NewComponentFromValue builds a component by reading the token and the text out
// of one string. A string with no deliminator is a token on its own.
func NewComponentFromValue(value string) *Component {
	token, text, found := strings.Cut(value, componentDeliminator)
	if !found {
		return NewComponent(value, "")
	}

	return NewComponent(token, text)
}

// GetToken returns the token of the component.
func (c *Component) GetToken() string {
	return c.token
}

// WithToken returns a copy of the component for another token.
func (c *Component) WithToken(token string) contract.ComponentContract {
	copied := *c
	copied.token = strings.TrimSpace(token)

	return &copied
}

// GetText returns the text of the component.
func (c *Component) GetText() string {
	return c.text
}

// WithText returns a copy of the component for another text.
func (c *Component) WithText(text string) contract.ComponentContract {
	copied := *c
	copied.text = strings.TrimSpace(text)

	return &copied
}

// String returns the whole component as a string. A component with no text is
// its token on its own.
func (c *Component) String() string {
	if c.token != "" && c.text != "" {
		return c.token + componentDeliminator + c.text
	}

	return c.token
}

// Value is one value of one header, built from components.
type Value struct {
	components []contract.ComponentContract
}

// NewValue builds a value from its components. It drops a component that renders
// to nothing, so a value never carries an empty part.
func NewValue(components ...contract.ComponentContract) *Value {
	return &Value{components: filterComponents(components)}
}

// NewValueFromValue builds a value by reading each component out of one string.
func NewValueFromValue(value string) *Value {
	parts := strings.Split(value, valueDeliminator)

	components := make([]contract.ComponentContract, 0, len(parts))

	for _, part := range parts {
		components = append(components, NewComponentFromValue(strings.TrimSpace(part)))
	}

	return NewValue(components...)
}

// GetComponents returns each component of the value.
func (v *Value) GetComponents() []contract.ComponentContract {
	return v.components
}

// WithComponents returns a copy of the value with other components.
func (v *Value) WithComponents(components ...contract.ComponentContract) contract.ValueContract {
	copied := *v
	copied.components = filterComponents(components)

	return &copied
}

// WithAddedComponents returns a copy of the value with the components added
// after the ones it holds.
func (v *Value) WithAddedComponents(components ...contract.ComponentContract) contract.ValueContract {
	copied := *v
	copied.components = append(slicesClone(v.components), filterComponents(components)...)

	return &copied
}

// String returns the whole value as a string, with each component separated by
// the value deliminator.
//
// The value holds no component that renders to nothing, because every path that
// sets the components filters them, so this needs no guard of its own.
func (v *Value) String() string {
	parts := make([]string, 0, len(v.components))

	for _, component := range v.components {
		parts = append(parts, strings.TrimSpace(component.String()))
	}

	return strings.Join(parts, valueDeliminator+" ")
}

// filterComponents drops each component that renders to nothing.
func filterComponents(components []contract.ComponentContract) []contract.ComponentContract {
	filtered := make([]contract.ComponentContract, 0, len(components))

	for _, component := range components {
		if strings.TrimSpace(component.String()) == "" {
			continue
		}

		filtered = append(filtered, component)
	}

	return filtered
}

// slicesClone copies the slice, so a copy of a value never shares its backing
// array with the value that it came from.
func slicesClone(components []contract.ComponentContract) []contract.ComponentContract {
	cloned := make([]contract.ComponentContract, len(components))

	copy(cloned, components)

	return cloned
}
