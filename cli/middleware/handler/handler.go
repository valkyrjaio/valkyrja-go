/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package handler runs the CLI middleware of one stage, in order.
//
// Each handler holds the binding key of every middleware that it runs, and it
// resolves the middleware from the container as it reaches it. A middleware
// receives the handler, so the middleware decides whether the run continues.
package handler

import (
	containerconstant "github.com/valkyrjaio/valkyrja-go/v26/container/constant"
	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
)

// Handler holds the middleware of one stage and walks them in order.
//
// Warning: a handler appends each middleware and never dedupes. A middleware
// that is added twice runs twice. That is the developer's error, and the
// framework does not correct it, because the generated cache must match what the
// runtime collects.
type Handler struct {
	container  containercontract.ContainerContract
	middleware []string
	index      int
}

// NewHandler builds a handler over a container, with the binding key of each
// middleware that it runs.
func NewHandler(container containercontract.ContainerContract, middleware ...string) Handler {
	return Handler{
		container:  container,
		middleware: middleware,
	}
}

// Add appends each middleware, by binding key, after the ones the handler holds.
func (h *Handler) Add(middleware ...string) {
	h.middleware = append(h.middleware, middleware...)
}

// hasNext reports whether the handler holds a middleware that it has not run.
func (h *Handler) hasNext() bool {
	return h.index < len(h.middleware)
}

// getNext resolves the next middleware and advances the handler past it.
//
// It returns nil where the handler holds no next middleware, and where the
// container resolves nothing for the binding key. A middleware that the
// container cannot resolve is skipped rather than fatal, so one unregistered
// middleware does not end the run.
func (h *Handler) getNext() any {
	if !h.hasNext() {
		return nil
	}

	id := h.middleware[h.index]
	h.index++

	resolved, err := h.container.Get(id, nil, containerconstant.NewInstanceOrThrowException)
	if err != nil {
		return nil
	}

	return resolved
}
