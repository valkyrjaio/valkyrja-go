/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package contract holds every contract of the application component.
//
// The package sits above the other components' contract packages and imports
// them. That is why the Go port builds the contract layer of every component
// before it builds the application: `ApplicationContract` names the CLI and the
// HTTP route provider contracts, and Go resolves an import at package level, so
// those packages must compile first. The other ports have no such constraint,
// because a type-only import may be circular in TypeScript and in PHP.
package contract

import (
	clicontract "github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	eventcontract "github.com/valkyrjaio/valkyrja-go/v26/event/contract"
	httpcontract "github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	throwablecontract "github.com/valkyrjaio/valkyrja-go/v26/throwable/contract"
)

// ApplicationCallbackFunc runs against the application once the application
// publishes its providers.
type ApplicationCallbackFunc func(app ApplicationContract)

// ApplicationContract is the application itself.
type ApplicationContract interface {
	// GetContainer returns the container of the application.
	GetContainer() containercontract.ContainerContract

	// PublishProviderCallbacks runs each callback that the config holds.
	PublishProviderCallbacks()

	// GetProviders returns every component provider, including the ones that a
	// component provider names.
	GetProviders() []ComponentProviderContract

	// GetContainerProviders returns every service provider of every component.
	GetContainerProviders() []containercontract.ServiceProviderContract

	// GetEventProviders returns every listener provider of every component.
	GetEventProviders() []eventcontract.ListenerProviderContract

	// GetCliProviders returns every CLI route provider of every component.
	GetCliProviders() []clicontract.CliRouteProviderContract

	// GetHttpProviders returns every HTTP route provider of every component.
	GetHttpProviders() []httpcontract.HttpRouteProviderContract

	// GetDebugMode reports whether the application runs in debug mode.
	GetDebugMode() bool

	// GetEnvironment returns the environment that the application runs in.
	GetEnvironment() string

	// GetVersion returns the version of the application.
	GetVersion() string
}

// ComponentProviderContract is the top-level provider of one component. It names
// the providers of every kind that the component registers.
//
// A provider returns a literal slice, and never a computed one, because `sindri`
// reads the slice from the source rather than by running it.
type ComponentProviderContract interface {
	// GetComponentProviders returns each component that this component needs.
	GetComponentProviders(app ApplicationContract) []ComponentProviderContract

	// GetContainerProviders returns each service provider of the component.
	GetContainerProviders(app ApplicationContract) []containercontract.ServiceProviderContract

	// GetEventProviders returns each listener provider of the component.
	GetEventProviders(app ApplicationContract) []eventcontract.ListenerProviderContract

	// GetCliProviders returns each CLI route provider of the component.
	GetCliProviders(app ApplicationContract) []clicontract.CliRouteProviderContract

	// GetHttpProviders returns each HTTP route provider of the component.
	GetHttpProviders(app ApplicationContract) []httpcontract.HttpRouteProviderContract
}

// PublishableComponentProviderContract is a component provider that runs
// something of its own once the application publishes it.
type PublishableComponentProviderContract interface {
	// Publish runs what the component needs at boot.
	Publish(app ApplicationContract)
}

// ConfigContract is the configuration of the application.
//
// The other ports declare each of these as a readonly property. Go has no
// property, so each one is a getter, which the method naming rules require in
// any case.
//
//nolint:interfacebloat // Parity with the PHP reference implementation.
type ConfigContract interface {
	// GetNamespace returns the namespace of the application.
	GetNamespace() string

	// GetDir returns the directory that the application runs from.
	GetDir() string

	// GetVersion returns the version of the application.
	GetVersion() string

	// GetEnvironment returns the environment that the application runs in.
	GetEnvironment() string

	// GetDebugMode reports whether the application runs in debug mode.
	GetDebugMode() bool

	// GetTimezone returns the timezone that the application runs in.
	GetTimezone() string

	// GetKey returns the secret key of the application.
	GetKey() string

	// GetDataPath returns the path that `sindri` writes the generated data to.
	GetDataPath() string

	// GetDataNamespace returns the namespace that `sindri` writes the generated
	// data under.
	GetDataNamespace() string

	// GetProviders returns each component provider that the application
	// registers.
	GetProviders() []ComponentProviderContract

	// GetCallbacks returns each callback that runs once the application
	// publishes its providers.
	GetCallbacks() []ApplicationCallbackFunc
}

// CliConfigContract is the configuration that a CLI application adds.
//
// Each middleware list holds a binding key, which is the framework's own class
// reference in every port that has no class token.
type CliConfigContract interface {
	ConfigContract

	// GetApplicationName returns the name that the CLI prints for itself.
	GetApplicationName() string

	// GetDefaultCommandName returns the command that runs where the caller
	// names none.
	GetDefaultCommandName() string

	// GetInputReceivedMiddleware returns the binding key of each
	// input-received middleware.
	GetInputReceivedMiddleware() []string

	// GetRouteMatchedMiddleware returns the binding key of each route-matched
	// middleware.
	GetRouteMatchedMiddleware() []string

	// GetRouteNotMatchedMiddleware returns the binding key of each
	// route-not-matched middleware.
	GetRouteNotMatchedMiddleware() []string

	// GetRouteDispatchedMiddleware returns the binding key of each
	// route-dispatched middleware.
	GetRouteDispatchedMiddleware() []string

	// GetThrowableCaughtMiddleware returns the binding key of each
	// throwable-caught middleware.
	GetThrowableCaughtMiddleware() []string

	// GetProcessExitingMiddleware returns the binding key of each
	// process-exiting middleware.
	GetProcessExitingMiddleware() []string
}

// HttpConfigContract is the configuration that an HTTP application adds.
//
// Each middleware list holds a binding key, which is the framework's own class
// reference in every port that has no class token.
type HttpConfigContract interface {
	ConfigContract

	// GetRequestReceivedMiddleware returns the binding key of each
	// request-received middleware.
	GetRequestReceivedMiddleware() []string

	// GetRouteMatchedMiddleware returns the binding key of each route-matched
	// middleware.
	GetRouteMatchedMiddleware() []string

	// GetRouteNotMatchedMiddleware returns the binding key of each
	// route-not-matched middleware.
	GetRouteNotMatchedMiddleware() []string

	// GetRouteDispatchedMiddleware returns the binding key of each
	// route-dispatched middleware.
	GetRouteDispatchedMiddleware() []string

	// GetThrowableCaughtMiddleware returns the binding key of each
	// throwable-caught middleware.
	GetThrowableCaughtMiddleware() []string

	// GetSendingResponseMiddleware returns the binding key of each
	// sending-response middleware.
	GetSendingResponseMiddleware() []string

	// GetResponseSentMiddleware returns the binding key of each response-sent
	// middleware.
	GetResponseSentMiddleware() []string
}

// ApplicationThrowable is the contract that every error of the application
// component satisfies.
type ApplicationThrowable interface {
	throwablecontract.ValkyrjaThrowable

	// IsApplicationThrowable marks the error as one that the application
	// raised. The mark is what separates this contract from the root contract,
	// which Go otherwise treats as the same type.
	IsApplicationThrowable() bool
}
