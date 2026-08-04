/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package processor_test

import (
	"regexp"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/data"
	"github.com/valkyrjaio/valkyrja-go/v26/http/routing/processor"
)

const (
	routeName = "users.show"
	usersPath = "/users"
)

// process runs the processor over the route and returns the result.
func process(route contract.RouteContract) contract.RouteContract {
	return processor.NewProcessor().Route(route)
}

// asDynamic reads the result as a dynamic route, and fails the test where it is
// not one.
func asDynamic(t *testing.T, route contract.RouteContract) contract.DynamicRouteContract {
	t.Helper()

	dynamic, isDynamic := route.(contract.DynamicRouteContract)
	if !isDynamic {
		t.Fatalf("the route must be dynamic, but is: %T", route)
	}

	return dynamic
}

func TestThePathTakesOneLeadingSeparatorAndNoTrailingOne(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"users":     usersPath,
		usersPath:   usersPath,
		"/users/":   usersPath,
		"//users//": usersPath,
		"":          "/",
	}

	for path, want := range tests {
		got := process(data.NewRoute(path, routeName, nil)).GetPath()

		if got != want {
			t.Errorf("the path %q must become %q, but became %q", path, want, got)
		}
	}
}

func TestAStaticRouteTakesNoRegex(t *testing.T) {
	t.Parallel()

	route := process(data.NewRoute(usersPath, routeName, nil))

	if asDynamic(t, route).GetRegex() != "" {
		t.Error("a static route must take no regular expression, but took one")
	}
}

func TestARouteWithNoParameterInItsPathTakesNoRegex(t *testing.T) {
	t.Parallel()

	route := process(data.NewDynamicRoute(
		data.NewRoute(usersPath, routeName, nil),
		"",
		data.NewParameter("id", constant.RegexNum),
	))

	if asDynamic(t, route).GetRegex() != "" {
		t.Error("a path with no parameter must take no regular expression, but took one")
	}
}

func TestTheProcessorBuildsAnAnchoredRegexFromThePath(t *testing.T) {
	t.Parallel()

	route := process(data.NewDynamicRoute(
		data.NewRoute("/users/{id}", routeName, nil),
		"",
		data.NewParameter("id", constant.RegexNum),
	))

	regex := asDynamic(t, route).GetRegex()

	if regex != `^\/users\/(?P<id>\d+)$` {
		t.Fatalf("the regular expression must be anchored and named, but is: %q", regex)
	}

	compiled, err := regexp.Compile(regex)
	if err != nil {
		t.Fatalf("the regular expression must compile under RE2, but reported: %v", err)
	}

	if !compiled.MatchString("/users/42") {
		t.Error("the regular expression must match the path, but did not")
	}

	if compiled.MatchString("/users/abc") {
		t.Error("the regular expression must reject a path that the parameter does not match, but did not")
	}

	if compiled.MatchString("/users/42/edit") {
		t.Error("the anchors must make the match exact, but a longer path matched")
	}

	matches := compiled.FindStringSubmatch("/users/42")

	if matches[compiled.SubexpIndex("id")] != "42" {
		t.Error("the parameter must read back by name, but did not")
	}
}

func TestTheProcessorKeepsARegexThatTheRouteCarries(t *testing.T) {
	t.Parallel()

	route := process(data.NewDynamicRoute(
		data.NewRoute("/users/{id}", routeName, nil),
		"^/kept$",
		data.NewParameter("id", constant.RegexNum),
	))

	if asDynamic(t, route).GetRegex() != "^/kept$" {
		t.Error("the processor must keep a regular expression that the route carries, but did not")
	}
}

func TestAnOptionalParameterTakesAnOptionalGroup(t *testing.T) {
	t.Parallel()

	route := process(data.NewDynamicRoute(
		data.NewRoute("/users{id?}", routeName, nil),
		"",
		data.NewParameter("id", constant.RegexNum),
	))

	regex := asDynamic(t, route).GetRegex()

	compiled, err := regexp.Compile(regex)
	if err != nil {
		t.Fatalf("the regular expression must compile under RE2, but reported: %v", err)
	}

	if !compiled.MatchString(usersPath) {
		t.Errorf("an optional parameter must let the path leave it out, but %q did not match /users", regex)
	}

	if regex != `^\/users(?:\/)?(?P<id>\d+)?$` {
		t.Errorf("the optional group must carry the separator, but the regular expression is: %q", regex)
	}

	if !compiled.MatchString("/users/42") {
		t.Errorf("an optional parameter must still match a path that carries it, but %q did not", regex)
	}

	if !asDynamic(t, route).GetParameters()[0].IsOptional() {
		t.Error("the parameter must be marked optional, but was not")
	}
}

func TestAParameterMarkedOptionalTakesAnOptionalGroup(t *testing.T) {
	t.Parallel()

	route := process(data.NewDynamicRoute(
		data.NewRoute("/users{id?}", routeName, nil),
		"",
		data.NewParameter("id", constant.RegexNum).WithIsOptional(true),
	))

	compiled, err := regexp.Compile(asDynamic(t, route).GetRegex())
	if err != nil {
		t.Fatalf("the regular expression must compile under RE2, but reported: %v", err)
	}

	if !compiled.MatchString(usersPath) {
		t.Error("a parameter that is marked optional must let the path leave it out, but did not")
	}
}

func TestAParameterThatDoesNotCaptureTakesANonCaptureGroup(t *testing.T) {
	t.Parallel()

	route := process(data.NewDynamicRoute(
		data.NewRoute("/users/{id}", routeName, nil),
		"",
		data.NewParameter("id", constant.RegexNum).WithShouldCapture(false),
	))

	regex := asDynamic(t, route).GetRegex()

	if regex != `^\/users\/(?:\d+)$` {
		t.Errorf("the parameter must take a non-capture group, but the regular expression is: %q", regex)
	}
}

func TestAParameterThatThePathDoesNotNameLeavesTheRegexAlone(t *testing.T) {
	t.Parallel()

	route := process(data.NewDynamicRoute(
		data.NewRoute("/users/{id}", routeName, nil),
		"",
		data.NewParameter("slug", constant.RegexSlug),
	))

	regex := asDynamic(t, route).GetRegex()

	if regex != `^\/users\/{id}$` {
		t.Errorf("a parameter that the path does not name must leave the path alone, but the regex is: %q", regex)
	}
}

func TestTheProcessorSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.ProcessorContract = processor.NewProcessor()

	if built.Route(data.NewRoute("users", routeName, nil)).GetPath() != usersPath {
		t.Error("the contract must prepare the route, but did not")
	}
}

func TestAnOptionalParameterNeedsItsPathToLeaveOutTheSeparator(t *testing.T) {
	t.Parallel()

	// The optional group carries the separator, so a path that also carries one
	// needs the separator and no longer matches without the parameter.
	route := process(data.NewDynamicRoute(
		data.NewRoute("/users/{id?}", routeName, nil),
		"",
		data.NewParameter("id", constant.RegexNum),
	))

	compiled, err := regexp.Compile(asDynamic(t, route).GetRegex())
	if err != nil {
		t.Fatalf("the regular expression must compile under RE2, but reported: %v", err)
	}

	if compiled.MatchString(usersPath) {
		t.Error("a path that carries the separator must not match without the parameter, but did")
	}

	if !compiled.MatchString("/users/42") {
		t.Error("the path must still match with the parameter, but did not")
	}
}

// staticOnlyRouteFixture is a route from outside the routing package: it carries
// a parameter in its path and satisfies no dynamic contract, so the processor
// returns it as it received it.
type staticOnlyRouteFixture struct {
	contract.RouteContract

	path string
}

// GetPath returns the path of the route.
func (r *staticOnlyRouteFixture) GetPath() string { return r.path }

// WithPath returns a copy of the route for another path.
func (r *staticOnlyRouteFixture) WithPath(path string) contract.RouteContract {
	return &staticOnlyRouteFixture{path: path}
}

func TestARouteThatSatisfiesNoDynamicContractIsReturnedAsItIs(t *testing.T) {
	t.Parallel()

	route := process(&staticOnlyRouteFixture{path: "/users/{id}"})

	if route.GetPath() != "/users/{id}" {
		t.Errorf("the processor must return the route as it is, but the path is: %q", route.GetPath())
	}

	if _, isDynamic := route.(contract.DynamicRouteContract); isDynamic {
		t.Error("the route must satisfy no dynamic contract, but did")
	}
}
