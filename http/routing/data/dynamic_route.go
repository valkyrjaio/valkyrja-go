/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data

import (
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
)

type Parameter struct {
	name          string
	regex         string
	cast          contract.CastFunc
	isOptional    bool
	shouldCapture bool
	defaultValue  any
	value         any
}

// NewParameter builds a parameter under a name, that matches a regular
// expression.
func NewParameter(name string, regex string) *Parameter {
	return &Parameter{
		name:          name,
		regex:         regex,
		shouldCapture: true,
	}
}

// GetName returns the name of the parameter.
func (p *Parameter) GetName() string {
	return p.name
}

// WithName returns a copy of the parameter under another name.
func (p *Parameter) WithName(name string) contract.ParameterContract {
	copied := *p
	copied.name = name

	return &copied
}

// GetRegex returns the regular expression that the parameter matches.
func (p *Parameter) GetRegex() string {
	return p.regex
}

// WithRegex returns a copy of the parameter for another regular expression.
func (p *Parameter) WithRegex(regex string) contract.ParameterContract {
	copied := *p
	copied.regex = regex

	return &copied
}

// HasCast reports whether the parameter casts its value to a type.
func (p *Parameter) HasCast() bool {
	return p.cast != nil
}

// GetCast returns what converts the value of the parameter, and nil where the
// parameter casts nothing.
func (p *Parameter) GetCast() contract.CastFunc {
	return p.cast
}

// WithCast returns a copy of the parameter for another cast.
func (p *Parameter) WithCast(cast contract.CastFunc) contract.ParameterContract {
	copied := *p
	copied.cast = cast

	return &copied
}

// IsOptional reports whether the path matches without the parameter.
func (p *Parameter) IsOptional() bool {
	return p.isOptional
}

// WithIsOptional returns a copy of the parameter with another optional flag.
func (p *Parameter) WithIsOptional(isOptional bool) contract.ParameterContract {
	copied := *p
	copied.isOptional = isOptional

	return &copied
}

// ShouldCapture reports whether the router passes the parameter to the handler.
func (p *Parameter) ShouldCapture() bool {
	return p.shouldCapture
}

// WithShouldCapture returns a copy of the parameter with another capture flag.
func (p *Parameter) WithShouldCapture(shouldCapture bool) contract.ParameterContract {
	copied := *p
	copied.shouldCapture = shouldCapture

	return &copied
}

// GetDefault returns the value that the router uses where the path carries none.
func (p *Parameter) GetDefault() any {
	return p.defaultValue
}

// WithDefault returns a copy of the parameter for another default value.
func (p *Parameter) WithDefault(defaultValue any) contract.ParameterContract {
	copied := *p
	copied.defaultValue = defaultValue

	return &copied
}

// GetValue returns the value that the path carried.
func (p *Parameter) GetValue() any {
	return p.value
}

// WithValue returns a copy of the parameter for another value.
func (p *Parameter) WithValue(value any) contract.ParameterContract {
	copied := *p
	copied.value = value

	return &copied
}

// NewDynamicRoute builds a route whose path carries a parameter, with the
// regular expression that its whole path matches.
func NewDynamicRoute(route *Route, regex string, parameters ...contract.ParameterContract) *Route {
	copied := *route
	copied.regex = regex
	copied.parameters = parameters

	return &copied
}

// IsDynamic reports whether the route carries a parameter.
func (r *Route) IsDynamic() bool {
	return r.regex != "" || len(r.parameters) > 0
}

// GetRegex returns the regular expression that the whole path matches.
func (r *Route) GetRegex() string {
	return r.regex
}

// WithRegex returns a copy of the route for another regular expression.
func (r *Route) WithRegex(regex string) contract.DynamicRouteContract {
	copied := *r
	copied.regex = regex

	return &copied
}

// GetParameters returns each parameter of the path.
func (r *Route) GetParameters() []contract.ParameterContract {
	return r.parameters
}

// WithParameters returns a copy of the route with other parameters.
func (r *Route) WithParameters(parameters ...contract.ParameterContract) contract.DynamicRouteContract {
	copied := *r
	copied.parameters = parameters

	return &copied
}

// WithAddedParameters returns a copy of the route with the parameters appended.
func (r *Route) WithAddedParameters(
	parameters ...contract.ParameterContract,
) contract.DynamicRouteContract {
	combined := make([]contract.ParameterContract, 0, len(r.parameters)+len(parameters))
	combined = append(combined, r.parameters...)
	combined = append(combined, parameters...)

	copied := *r
	copied.parameters = combined

	return &copied
}
