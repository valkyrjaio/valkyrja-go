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

type HttpConfig struct {
	Config

	RequestReceivedMiddleware []string
	RouteMatchedMiddleware    []string
	RouteNotMatchedMiddleware []string
	RouteDispatchedMiddleware []string
	ThrowableCaughtMiddleware []string
	SendingResponseMiddleware []string
	ResponseSentMiddleware    []string
}

// NewHttpConfig builds the configuration that every field takes its default
// value in.
func NewHttpConfig(providers ...contract.ComponentProviderContract) *HttpConfig {
	return &HttpConfig{
		Config:                    *NewConfig(providers...),
		RequestReceivedMiddleware: []string{},
		RouteMatchedMiddleware:    []string{},
		RouteNotMatchedMiddleware: []string{},
		RouteDispatchedMiddleware: []string{},
		ThrowableCaughtMiddleware: []string{},
		SendingResponseMiddleware: []string{},
		ResponseSentMiddleware:    []string{},
	}
}

// GetRequestReceivedMiddleware returns the binding key of each request-received
// middleware.
func (c *HttpConfig) GetRequestReceivedMiddleware() []string {
	return c.RequestReceivedMiddleware
}

// GetRouteMatchedMiddleware returns the binding key of each route-matched
// middleware.
func (c *HttpConfig) GetRouteMatchedMiddleware() []string {
	return c.RouteMatchedMiddleware
}

// GetRouteNotMatchedMiddleware returns the binding key of each route-not-matched
// middleware.
func (c *HttpConfig) GetRouteNotMatchedMiddleware() []string {
	return c.RouteNotMatchedMiddleware
}

// GetRouteDispatchedMiddleware returns the binding key of each route-dispatched
// middleware.
func (c *HttpConfig) GetRouteDispatchedMiddleware() []string {
	return c.RouteDispatchedMiddleware
}

// GetThrowableCaughtMiddleware returns the binding key of each throwable-caught
// middleware.
func (c *HttpConfig) GetThrowableCaughtMiddleware() []string {
	return c.ThrowableCaughtMiddleware
}

// GetSendingResponseMiddleware returns the binding key of each sending-response
// middleware.
func (c *HttpConfig) GetSendingResponseMiddleware() []string {
	return c.SendingResponseMiddleware
}

// GetResponseSentMiddleware returns the binding key of each response-sent
// middleware.
func (c *HttpConfig) GetResponseSentMiddleware() []string {
	return c.ResponseSentMiddleware
}
