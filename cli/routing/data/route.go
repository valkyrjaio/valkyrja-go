/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
)

type Route struct {
	name        string
	description string
	handler     contract.CliHandlerFunc
	helpText    contract.HelpTextFunc

	arguments []contract.ArgumentParameterContract
	options   []contract.OptionParameterContract

	routeMatchedMiddleware    []string
	routeDispatchedMiddleware []string
	throwableCaughtMiddleware []string
	processExitingMiddleware  []string
}

// NewRoute builds a command that runs the handler.
func NewRoute(name string, description string, handler contract.CliHandlerFunc) *Route {
	return &Route{
		name:                      name,
		description:               description,
		handler:                   handler,
		arguments:                 []contract.ArgumentParameterContract{},
		options:                   []contract.OptionParameterContract{},
		routeMatchedMiddleware:    []string{},
		routeDispatchedMiddleware: []string{},
		throwableCaughtMiddleware: []string{},
		processExitingMiddleware:  []string{},
	}
}

// GetName returns the name of the command.
func (r *Route) GetName() string {
	return r.name
}

// WithName returns a copy of the route under another name.
func (r *Route) WithName(name string) contract.RouteContract {
	copied := *r
	copied.name = name

	return &copied
}

// GetDescription returns the description that the help text prints.
func (r *Route) GetDescription() string {
	return r.description
}

// WithDescription returns a copy of the route with another description.
func (r *Route) WithDescription(description string) contract.RouteContract {
	copied := *r
	copied.description = description

	return &copied
}

// HasHelpText reports whether the route builds its own help text.
func (r *Route) HasHelpText() bool {
	return r.helpText != nil
}

// GetHelpText returns what builds the help text, and nil where the route builds
// none.
func (r *Route) GetHelpText() contract.HelpTextFunc {
	return r.helpText
}

// GetHelpTextMessage returns the help text as a message, and nil where the route
// builds none.
func (r *Route) GetHelpTextMessage() contract.MessageContract {
	if r.helpText == nil {
		return nil
	}

	return r.helpText()
}

// WithHelpText returns a copy of the route that builds another help text.
func (r *Route) WithHelpText(helpText contract.HelpTextFunc) contract.RouteContract {
	copied := *r
	copied.helpText = helpText

	return &copied
}

// HasArguments reports whether the command takes an argument.
func (r *Route) HasArguments() bool {
	return len(r.arguments) > 0
}

// GetArguments returns each argument of the command.
func (r *Route) GetArguments() []contract.ArgumentParameterContract {
	return r.arguments
}

// HasArgument reports whether the command takes the argument.
func (r *Route) HasArgument(name string) bool {
	return r.GetArgument(name) != nil
}

// GetArgument returns the argument under the name, and nil where the command
// takes no argument under it.
func (r *Route) GetArgument(name string) contract.ArgumentParameterContract {
	for _, argument := range r.arguments {
		if argument.GetName() == name {
			return argument
		}
	}

	return nil
}

// WithArguments returns a copy of the route with other arguments.
func (r *Route) WithArguments(arguments ...contract.ArgumentParameterContract) contract.RouteContract {
	copied := *r
	copied.arguments = arguments

	return &copied
}

// WithAddedArguments returns a copy of the route with the arguments appended.
func (r *Route) WithAddedArguments(arguments ...contract.ArgumentParameterContract) contract.RouteContract {
	combined := make([]contract.ArgumentParameterContract, 0, len(r.arguments)+len(arguments))
	combined = append(combined, r.arguments...)
	combined = append(combined, arguments...)

	copied := *r
	copied.arguments = combined

	return &copied
}

// HasOptions reports whether the command takes an option.
func (r *Route) HasOptions() bool {
	return len(r.options) > 0
}

// GetOptions returns each option of the command.
func (r *Route) GetOptions() []contract.OptionParameterContract {
	return r.options
}

// HasOption reports whether the command takes the option.
func (r *Route) HasOption(name string) bool {
	return r.GetOption(name) != nil
}

// GetOption returns the option under the name, and nil where the command takes
// no option under it.
func (r *Route) GetOption(name string) contract.OptionParameterContract {
	for _, option := range r.options {
		if option.GetName() == name {
			return option
		}
	}

	return nil
}

// WithOptions returns a copy of the route with other options.
func (r *Route) WithOptions(options ...contract.OptionParameterContract) contract.RouteContract {
	copied := *r
	copied.options = options

	return &copied
}

// WithAddedOptions returns a copy of the route with the options appended.
func (r *Route) WithAddedOptions(options ...contract.OptionParameterContract) contract.RouteContract {
	combined := make([]contract.OptionParameterContract, 0, len(r.options)+len(options))
	combined = append(combined, r.options...)
	combined = append(combined, options...)

	copied := *r
	copied.options = combined

	return &copied
}

// GetRouteMatchedMiddleware returns the binding key of each route-matched
// middleware.
func (r *Route) GetRouteMatchedMiddleware() []string {
	return r.routeMatchedMiddleware
}

// WithRouteMatchedMiddleware returns a copy of the route with other
// route-matched middleware.
func (r *Route) WithRouteMatchedMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.routeMatchedMiddleware = middleware

	return &copied
}

// WithAddedRouteMatchedMiddleware returns a copy of the route with the
// route-matched middleware appended.
func (r *Route) WithAddedRouteMatchedMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.routeMatchedMiddleware = appendMiddleware(r.routeMatchedMiddleware, middleware)

	return &copied
}

// GetRouteDispatchedMiddleware returns the binding key of each route-dispatched
// middleware.
func (r *Route) GetRouteDispatchedMiddleware() []string {
	return r.routeDispatchedMiddleware
}

// WithRouteDispatchedMiddleware returns a copy of the route with other
// route-dispatched middleware.
func (r *Route) WithRouteDispatchedMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.routeDispatchedMiddleware = middleware

	return &copied
}

// WithAddedRouteDispatchedMiddleware returns a copy of the route with the
// route-dispatched middleware appended.
func (r *Route) WithAddedRouteDispatchedMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.routeDispatchedMiddleware = appendMiddleware(r.routeDispatchedMiddleware, middleware)

	return &copied
}

// GetThrowableCaughtMiddleware returns the binding key of each throwable-caught
// middleware.
func (r *Route) GetThrowableCaughtMiddleware() []string {
	return r.throwableCaughtMiddleware
}

// WithThrowableCaughtMiddleware returns a copy of the route with other
// throwable-caught middleware.
func (r *Route) WithThrowableCaughtMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.throwableCaughtMiddleware = middleware

	return &copied
}

// WithAddedThrowableCaughtMiddleware returns a copy of the route with the
// throwable-caught middleware appended.
func (r *Route) WithAddedThrowableCaughtMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.throwableCaughtMiddleware = appendMiddleware(r.throwableCaughtMiddleware, middleware)

	return &copied
}

// GetProcessExitingMiddleware returns the binding key of each process-exiting
// middleware.
func (r *Route) GetProcessExitingMiddleware() []string {
	return r.processExitingMiddleware
}

// WithProcessExitingMiddleware returns a copy of the route with other
// process-exiting middleware.
func (r *Route) WithProcessExitingMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.processExitingMiddleware = middleware

	return &copied
}

// WithAddedProcessExitingMiddleware returns a copy of the route with the
// process-exiting middleware appended.
func (r *Route) WithAddedProcessExitingMiddleware(middleware ...string) contract.RouteContract {
	copied := *r
	copied.processExitingMiddleware = appendMiddleware(r.processExitingMiddleware, middleware)

	return &copied
}

// GetHandler returns what the command runs.
func (r *Route) GetHandler() contract.CliHandlerFunc {
	return r.handler
}

// WithHandler returns a copy of the route that runs another handler.
func (r *Route) WithHandler(handler contract.CliHandlerFunc) contract.RouteContract {
	copied := *r
	copied.handler = handler

	return &copied
}

// appendMiddleware returns the middleware with the added middleware after it, in
// a slice of its own.
//
// Warning: the framework never dedupes middleware. A middleware that is added
// twice runs twice, because the generated cache must match what the runtime
// collects.
func appendMiddleware(middleware []string, added []string) []string {
	combined := make([]string, 0, len(middleware)+len(added))
	combined = append(combined, middleware...)
	combined = append(combined, added...)

	return combined
}
