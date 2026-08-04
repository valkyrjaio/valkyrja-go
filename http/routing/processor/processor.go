/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package processor prepares a route before the collection holds it.
package processor

import (
	"strings"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/constant"
)

// parameterStart opens a parameter in a route path.
const parameterStart = "{"

// parameterEnd closes a parameter in a route path.
const parameterEnd = "}"

// optionalMarker marks a parameter that the path may leave out.
const optionalMarker = "?"

// Processor prepares a route before the collection holds it.
//
// It normalizes the path, and it builds the regular expression of a dynamic
// route that carries none.
type Processor struct{}

// NewProcessor builds the processor.
func NewProcessor() *Processor {
	return &Processor{}
}

// Route returns the route, prepared.
//
// The path takes one leading separator and no trailing one, so two routes that a
// developer wrote differently file under the same path.
//
// Warning: an optional parameter carries its own separator, so its path leaves
// the separator out. Write `/users{id?}`, never `/users/{id?}` — the second one
// needs the separator, so `/users` does not match it.
func (p *Processor) Route(route contract.RouteContract) contract.RouteContract {
	route = route.WithPath("/" + strings.Trim(route.GetPath(), "/"))

	if !strings.Contains(route.GetPath(), parameterStart) {
		return route
	}

	// Every route satisfies the dynamic contract, because Go compares an
	// interface by its method set and one struct carries both shapes. A route
	// that does not is a type from outside this package, and it has no regular
	// expression to build, so the path guard above already returned it.
	dynamic, isDynamic := route.(contract.DynamicRouteContract)
	if !isDynamic {
		return route
	}

	return p.withRegexFromPath(dynamic)
}

// withRegexFromPath builds the regular expression of the route from its path. It
// returns the route unchanged where the route carries one already.
func (p *Processor) withRegexFromPath(route contract.DynamicRouteContract) contract.RouteContract {
	if route.GetRegex() != "" {
		return route
	}

	regex := strings.ReplaceAll(route.GetPath(), "/", constant.RegexPath)

	parameters := route.GetParameters()

	processed := make([]contract.ParameterContract, 0, len(parameters))

	for _, parameter := range parameters {
		prepared := p.withOptionalFromRegex(parameter, regex)

		replaced, found := p.replaceParameterName(prepared, regex)
		if !found {
			// A parameter that the path does not name cannot be filled, so the
			// route keeps the parameter and the regular expression stays as the
			// path left it. The matcher then gives the parameter its default.
			processed = append(processed, prepared)

			continue
		}

		regex = replaced

		processed = append(processed, prepared)
	}

	return route.
		WithParameters(processed...).
		WithRegex(constant.RegexStart + regex + constant.RegexEnd)
}

// withOptionalFromRegex returns the parameter marked optional where the path
// marks it so.
func (p *Processor) withOptionalFromRegex(
	parameter contract.ParameterContract,
	regex string,
) contract.ParameterContract {
	if parameter.IsOptional() {
		return parameter
	}

	if strings.Contains(regex, parameter.GetName()+optionalMarker) {
		return parameter.WithIsOptional(true)
	}

	return parameter
}

// replaceParameterName replaces the name of the parameter in the regular
// expression with the group that matches it. It reports whether the regular
// expression names the parameter.
func (p *Processor) replaceParameterName(
	parameter contract.ParameterContract,
	regex string,
) (string, bool) {
	name := p.getNameReplacement(parameter)

	if !strings.Contains(regex, name) {
		return regex, false
	}

	return strings.Replace(regex, name, p.getParameterRegex(parameter), 1), true
}

// getNameReplacement returns the parameter as the path spells it.
func (p *Processor) getNameReplacement(parameter contract.ParameterContract) string {
	if parameter.IsOptional() {
		return parameterStart + parameter.GetName() + optionalMarker + parameterEnd
	}

	return parameterStart + parameter.GetName() + parameterEnd
}

// getParameterRegex returns the group that matches the parameter.
func (p *Processor) getParameterRegex(parameter contract.ParameterContract) string {
	built := &strings.Builder{}

	if parameter.IsOptional() {
		built.WriteString(constant.RegexStartOptionalCaptureGroup)
	}

	if parameter.ShouldCapture() {
		built.WriteString(constant.RegexStartCaptureGroup)
		built.WriteString(constant.RegexStartCaptureGroupName)
		built.WriteString(parameter.GetName())
		built.WriteString(constant.RegexEndCaptureGroupName)
	} else {
		built.WriteString(constant.RegexStartNonCaptureGroup)
	}

	built.WriteString(parameter.GetRegex())

	if parameter.IsOptional() {
		built.WriteString(constant.RegexEndOptionalCaptureGroup)
	} else {
		built.WriteString(constant.RegexEndCaptureGroup)
	}

	return built.String()
}
