/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package contract

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/routing/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

type CliHandlerFunc func(container containercontract.ContainerContract, route RouteContract) OutputContract

type CastFunc func(value string) (any, error)

type HelpTextFunc func() MessageContract

//nolint:interfacebloat // Parity with the PHP reference implementation.
type ParameterContract interface {
	// GetName returns the name of the parameter.
	GetName() string

	// WithName returns a copy of the parameter under another name.
	WithName(name string) ParameterContract

	// HasCast reports whether the parameter casts its values to a type.
	HasCast() bool

	// GetCast returns what converts the values of the parameter, and nil where
	// the parameter casts nothing.
	GetCast() CastFunc

	// WithCast returns a copy of the parameter for another cast.
	WithCast(cast CastFunc) ParameterContract

	// WithoutCast returns a copy of the parameter with no cast.
	WithoutCast() ParameterContract

	// GetDescription returns the description that the help text prints.
	GetDescription() string

	// WithDescription returns a copy of the parameter with another
	// description.
	WithDescription(description string) ParameterContract

	// GetCastValues returns each value, cast to the type that the parameter
	// names, and reports a failure where a value does not convert.
	GetCastValues() ([]any, error)

	// HasFirstValue reports whether the caller gave the parameter a value.
	HasFirstValue() bool

	// GetFirstValue returns the first value that the caller gave.
	GetFirstValue() string

	// AreValuesValid reports whether every value is one that the parameter
	// accepts.
	AreValuesValid() bool

	// ValidateValues reports a failure where a value is one that the parameter
	// does not accept.
	//
	// The other ports throw here and return the parameter, so a caller chains
	// the call. Go reports a failure with an error instead.
	ValidateValues() error
}

type ArgumentParameterContract interface {
	ParameterContract

	// GetMode returns whether the command needs the argument.
	GetMode() constant.ArgumentMode

	// WithMode returns a copy of the argument in another mode.
	WithMode(mode constant.ArgumentMode) ArgumentParameterContract

	// GetValueMode returns how many values the argument takes.
	GetValueMode() constant.ArgumentValueMode

	// WithValueMode returns a copy of the argument in another value mode.
	WithValueMode(valueMode constant.ArgumentValueMode) ArgumentParameterContract

	// GetArguments returns each argument that the caller typed for this
	// parameter.
	GetArguments() []ArgumentContract

	// WithArguments returns a copy of the parameter with other arguments.
	WithArguments(arguments ...ArgumentContract) ArgumentParameterContract

	// WithAddedArguments returns a copy of the parameter with the arguments
	// appended.
	WithAddedArguments(arguments ...ArgumentContract) ArgumentParameterContract
}

//nolint:interfacebloat // Parity with the PHP reference implementation.
type OptionParameterContract interface {
	ParameterContract

	// GetShortNames returns each one-letter name of the option.
	GetShortNames() []string

	// WithShortNames returns a copy of the option with other short names.
	WithShortNames(shortNames ...string) OptionParameterContract

	// WithAddedShortNames returns a copy of the option with the short names
	// appended.
	WithAddedShortNames(shortNames ...string) OptionParameterContract

	// GetMode returns whether the command needs the option.
	GetMode() constant.OptionMode

	// WithMode returns a copy of the option in another mode.
	WithMode(mode constant.OptionMode) OptionParameterContract

	// GetValueMode returns how many values the option takes.
	GetValueMode() constant.OptionValueMode

	// WithValueMode returns a copy of the option in another value mode.
	WithValueMode(valueMode constant.OptionValueMode) OptionParameterContract

	// HasValueDisplayName reports whether the help text names the value.
	HasValueDisplayName() bool

	// GetValueDisplayName returns the name that the help text gives the value.
	GetValueDisplayName() string

	// WithValueDisplayName returns a copy of the option with another display
	// name.
	WithValueDisplayName(valueName string) OptionParameterContract

	// GetValidValues returns each value that the option accepts.
	GetValidValues() []string

	// WithValidValues returns a copy of the option that accepts other values.
	WithValidValues(validValues ...string) OptionParameterContract

	// WithAddedValidValues returns a copy of the option with the valid values
	// appended.
	WithAddedValidValues(validValues ...string) OptionParameterContract

	// HasDefaultValue reports whether the option has a value that the caller
	// leaves out.
	HasDefaultValue() bool

	// GetDefaultValue returns the value that the option uses where the caller
	// gives none.
	GetDefaultValue() string

	// WithDefaultValue returns a copy of the option with another default value.
	WithDefaultValue(defaultValue string) OptionParameterContract

	// GetOptions returns each option that the caller typed for this parameter.
	GetOptions() []OptionContract

	// WithOptions returns a copy of the parameter with other options, and
	// reports a failure where an option carries a value that the parameter takes
	// no value for.
	WithOptions(options ...OptionContract) (OptionParameterContract, error)

	// WithAddedOptions returns a copy of the parameter with the options
	// appended, and reports a failure the way `WithOptions` does.
	WithAddedOptions(options ...OptionContract) (OptionParameterContract, error)
}

//nolint:interfacebloat // Parity with the PHP reference implementation.
type RouteContract interface {
	// GetName returns the name of the command.
	GetName() string

	// WithName returns a copy of the route under another name.
	WithName(name string) RouteContract

	// GetDescription returns the description that the help text prints.
	GetDescription() string

	// WithDescription returns a copy of the route with another description.
	WithDescription(description string) RouteContract

	// HasHelpText reports whether the route builds its own help text.
	HasHelpText() bool

	// GetHelpText returns what builds the help text.
	GetHelpText() HelpTextFunc

	// GetHelpTextMessage returns the help text as a message.
	GetHelpTextMessage() MessageContract

	// WithHelpText returns a copy of the route that builds another help text.
	WithHelpText(helpText HelpTextFunc) RouteContract

	// HasArguments reports whether the command takes an argument.
	HasArguments() bool

	// GetArguments returns each argument of the command.
	GetArguments() []ArgumentParameterContract

	// HasArgument reports whether the command takes the argument.
	HasArgument(name string) bool

	// GetArgument returns the argument under the name.
	GetArgument(name string) ArgumentParameterContract

	// WithArguments returns a copy of the route with other arguments.
	WithArguments(arguments ...ArgumentParameterContract) RouteContract

	// WithAddedArguments returns a copy of the route with the arguments
	// appended.
	WithAddedArguments(arguments ...ArgumentParameterContract) RouteContract

	// HasOptions reports whether the command takes an option.
	HasOptions() bool

	// GetOptions returns each option of the command.
	GetOptions() []OptionParameterContract

	// HasOption reports whether the command takes the option.
	HasOption(name string) bool

	// GetOption returns the option under the name.
	GetOption(name string) OptionParameterContract

	// WithOptions returns a copy of the route with other options.
	WithOptions(options ...OptionParameterContract) RouteContract

	// WithAddedOptions returns a copy of the route with the options appended.
	WithAddedOptions(options ...OptionParameterContract) RouteContract

	// GetRouteMatchedMiddleware returns the binding key of each route-matched
	// middleware.
	GetRouteMatchedMiddleware() []string

	// WithRouteMatchedMiddleware returns a copy of the route with other
	// route-matched middleware.
	WithRouteMatchedMiddleware(middleware ...string) RouteContract

	// WithAddedRouteMatchedMiddleware returns a copy of the route with the
	// route-matched middleware appended.
	WithAddedRouteMatchedMiddleware(middleware ...string) RouteContract

	// GetRouteDispatchedMiddleware returns the binding key of each
	// route-dispatched middleware.
	GetRouteDispatchedMiddleware() []string

	// WithRouteDispatchedMiddleware returns a copy of the route with other
	// route-dispatched middleware.
	WithRouteDispatchedMiddleware(middleware ...string) RouteContract

	// WithAddedRouteDispatchedMiddleware returns a copy of the route with the
	// route-dispatched middleware appended.
	WithAddedRouteDispatchedMiddleware(middleware ...string) RouteContract

	// GetThrowableCaughtMiddleware returns the binding key of each
	// throwable-caught middleware.
	GetThrowableCaughtMiddleware() []string

	// WithThrowableCaughtMiddleware returns a copy of the route with other
	// throwable-caught middleware.
	WithThrowableCaughtMiddleware(middleware ...string) RouteContract

	// WithAddedThrowableCaughtMiddleware returns a copy of the route with the
	// throwable-caught middleware appended.
	WithAddedThrowableCaughtMiddleware(middleware ...string) RouteContract

	// GetProcessExitingMiddleware returns the binding key of each
	// process-exiting middleware.
	GetProcessExitingMiddleware() []string

	// WithProcessExitingMiddleware returns a copy of the route with other
	// process-exiting middleware.
	WithProcessExitingMiddleware(middleware ...string) RouteContract

	// WithAddedProcessExitingMiddleware returns a copy of the route with the
	// process-exiting middleware appended.
	WithAddedProcessExitingMiddleware(middleware ...string) RouteContract

	// GetHandler returns what the command runs.
	GetHandler() CliHandlerFunc

	// WithHandler returns a copy of the route that runs another handler.
	WithHandler(handler CliHandlerFunc) RouteContract
}

type CliRoutingDataContract interface {
	// GetRoutes returns each route, keyed by its own name.
	GetRoutes() map[string]RouteContract
}

type RouteCollectionContract interface {
	// GetData returns the collection's state.
	GetData() CliRoutingDataContract

	// SetFromData replaces the collection's state.
	SetFromData(data CliRoutingDataContract)

	// Add files each command under its own name.
	Add(commands ...RouteContract) RouteCollectionContract

	// Get returns the command under the name.
	Get(name string) RouteContract

	// Has reports whether the collection holds a command under the name.
	Has(name string) bool

	// All returns every command, keyed by its own name.
	All() map[string]RouteContract
}

type RouterContract interface {
	// Dispatch matches the input to a command and runs it.
	Dispatch(input InputContract) OutputContract

	// DispatchRoute runs the command for the input.
	DispatchRoute(input InputContract, route RouteContract) OutputContract
}

type InputHandlerContract interface {
	// Handle returns the output for the input.
	Handle(input InputContract) OutputContract

	// Exit ends the process with the exit code that the output holds.
	Exit(input InputContract, output OutputContract)

	// Run handles the input and exits.
	Run(input InputContract)
}

type CliRouteProviderContract interface {
	// GetRoutes returns each command that the component registers.
	GetRoutes() []RouteContract
}
