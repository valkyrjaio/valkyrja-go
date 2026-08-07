/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package matcher finds the route that a request path matches.
package matcher

import (
	"regexp"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
)

type Matcher struct {
	collection contract.RouteCollectionContract
}

// NewMatcher builds a matcher over a route collection.
func NewMatcher(collection contract.RouteCollectionContract) *Matcher {
	return &Matcher{collection: collection}
}

// Match returns the route that the path matches, and nil where none matches.
func (m *Matcher) Match(path string, requestMethod constant.RequestMethod) contract.RouteContract {
	if route := m.MatchStatic(path, requestMethod); route != nil {
		return route
	}

	return m.MatchDynamic(path, requestMethod)
}

// MatchStatic returns the static route that the path matches, and nil where none
// matches.
func (m *Matcher) MatchStatic(path string, requestMethod constant.RequestMethod) contract.RouteContract {
	if !m.collection.HasPath(path, requestMethod) {
		return nil
	}

	return m.collection.GetByPath(path, requestMethod)
}

// MatchDynamic returns the dynamic route that the path matches, and nil where
// none matches.
func (m *Matcher) MatchDynamic(path string, requestMethod constant.RequestMethod) contract.RouteContract {
	for regex := range m.collection.GetRegexes(requestMethod) {
		compiled, err := regexp.Compile(regex)
		if err != nil {
			continue
		}

		matches := compiled.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		route := m.collection.GetByRegex(regex, requestMethod)
		if route == nil {
			continue
		}

		return withParameterValues(route, compiled, matches)
	}

	return nil
}

// withParameterValues returns the route with the value that the path filled into
// each parameter.
func withParameterValues(
	route contract.DynamicRouteContract,
	compiled *regexp.Regexp,
	matches []string,
) contract.DynamicRouteContract {
	parameters := route.GetParameters()

	filled := make([]contract.ParameterContract, 0, len(parameters))

	for _, parameter := range parameters {
		filled = append(filled, withValueFromMatches(parameter, compiled, matches))
	}

	return route.WithParameters(filled...)
}

// withValueFromMatches returns the parameter with the value that the path
// carried, and with its default value where the path carried none.
func withValueFromMatches(
	parameter contract.ParameterContract,
	compiled *regexp.Regexp,
	matches []string,
) contract.ParameterContract {
	index := compiled.SubexpIndex(parameter.GetName())

	if index < 0 || index >= len(matches) || matches[index] == "" {
		return parameter.WithValue(parameter.GetDefault())
	}

	return parameter.WithValue(castMatch(parameter, matches[index]))
}

// castMatch returns the value that the cast of the parameter returned.
//
// Warning: the regular expression is what states the shape of a value, and the
// cast converts what the regular expression accepted already. A cast that
// reports a failure therefore means the two disagree, which is the developer's
// error. The parameter then carries the text as the path held it, and the
// handler reads a string where it declared another type.
func castMatch(parameter contract.ParameterContract, match string) any {
	cast := parameter.GetCast()
	if cast == nil {
		return match
	}

	value, err := cast(match)
	if err != nil {
		return match
	}

	return value
}
