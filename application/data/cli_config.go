/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package data

import (
	"github.com/valkyrjaio/valkyrja-go/v26/application/contract"
)

// DefaultApplicationName is the name that the CLI prints for an application that
// states none.
const DefaultApplicationName = "Valkyrja"

// DefaultCommandName is the command that runs where the caller names none.
const DefaultCommandName = "list"

type CliConfig struct {
	Config

	ApplicationName    string
	DefaultCommandName string

	InputReceivedMiddleware   []string
	RouteMatchedMiddleware    []string
	RouteNotMatchedMiddleware []string
	RouteDispatchedMiddleware []string
	ThrowableCaughtMiddleware []string
	ProcessExitingMiddleware  []string
}

// NewCliConfig builds the configuration that every field takes its default value
// in.
func NewCliConfig(providers ...contract.ComponentProviderContract) *CliConfig {
	return &CliConfig{
		Config:                    *NewConfig(providers...),
		ApplicationName:           DefaultApplicationName,
		DefaultCommandName:        DefaultCommandName,
		InputReceivedMiddleware:   []string{},
		RouteMatchedMiddleware:    []string{},
		RouteNotMatchedMiddleware: []string{},
		RouteDispatchedMiddleware: []string{},
		ThrowableCaughtMiddleware: []string{},
		ProcessExitingMiddleware:  []string{},
	}
}

// GetApplicationName returns the name that the CLI prints for itself.
func (c *CliConfig) GetApplicationName() string {
	return c.ApplicationName
}

// GetDefaultCommandName returns the command that runs where the caller names
// none.
func (c *CliConfig) GetDefaultCommandName() string {
	return c.DefaultCommandName
}

// GetInputReceivedMiddleware returns the binding key of each input-received
// middleware.
func (c *CliConfig) GetInputReceivedMiddleware() []string {
	return c.InputReceivedMiddleware
}

// GetRouteMatchedMiddleware returns the binding key of each route-matched
// middleware.
func (c *CliConfig) GetRouteMatchedMiddleware() []string {
	return c.RouteMatchedMiddleware
}

// GetRouteNotMatchedMiddleware returns the binding key of each route-not-matched
// middleware.
func (c *CliConfig) GetRouteNotMatchedMiddleware() []string {
	return c.RouteNotMatchedMiddleware
}

// GetRouteDispatchedMiddleware returns the binding key of each route-dispatched
// middleware.
func (c *CliConfig) GetRouteDispatchedMiddleware() []string {
	return c.RouteDispatchedMiddleware
}

// GetThrowableCaughtMiddleware returns the binding key of each throwable-caught
// middleware.
func (c *CliConfig) GetThrowableCaughtMiddleware() []string {
	return c.ThrowableCaughtMiddleware
}

// GetProcessExitingMiddleware returns the binding key of each process-exiting
// middleware.
func (c *CliConfig) GetProcessExitingMiddleware() []string {
	return c.ProcessExitingMiddleware
}
