/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package handler is the server's entry point for one request.
package handler

import (
	"fmt"
	"net/http"

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	serverconstant "github.com/valkyrjaio/valkyrja-go/v26/http/server/constant"
)

type RequestHandler struct {
	container containercontract.ContainerContract
	router    contract.RouterContract

	requestReceivedHandler contract.RequestReceivedHandlerContract
	throwableCaughtHandler contract.ThrowableCaughtHandlerContract
	sendingResponseHandler contract.SendingResponseHandlerContract
	responseSentHandler    contract.ResponseSentHandlerContract

	responseFactory contract.ResponseFactoryContract
	debug           bool
}

// NewRequestHandler builds the handler over a container, a router, the
// middleware handler of each stage, and a response factory.
func NewRequestHandler(
	container containercontract.ContainerContract,
	router contract.RouterContract,
	requestReceivedHandler contract.RequestReceivedHandlerContract,
	throwableCaughtHandler contract.ThrowableCaughtHandlerContract,
	sendingResponseHandler contract.SendingResponseHandlerContract,
	responseSentHandler contract.ResponseSentHandlerContract,
	responseFactory contract.ResponseFactoryContract,
	debug bool,
) *RequestHandler {
	return &RequestHandler{
		container:              container,
		router:                 router,
		requestReceivedHandler: requestReceivedHandler,
		throwableCaughtHandler: throwableCaughtHandler,
		sendingResponseHandler: sendingResponseHandler,
		responseSentHandler:    responseSentHandler,
		responseFactory:        responseFactory,
		debug:                  debug,
	}
}

// Handle returns the response for the request.
func (h *RequestHandler) Handle(request contract.ServerRequestContract) contract.ResponseContract {
	response := h.dispatch(request)

	h.container.SetSingleton(serverconstant.ResponseContractServiceID, response)

	return response
}

// Send writes the response to the writer.
func (h *RequestHandler) Send(
	response contract.ResponseContract,
	writer http.ResponseWriter,
) contract.RequestHandlerContract {
	for _, header := range response.GetHeaders().GetAll() {
		writer.Header().Set(header.GetName(), header.GetHeaderLine())
	}

	writer.WriteHeader(int(response.GetStatusCode()))

	body := response.GetBody()

	// A body that no reader reads writes nothing, which is what an empty
	// response carries.
	_ = body.Rewind()

	contents, err := body.GetContents()
	if err != nil {
		return h
	}

	// The writer reports how many bytes it took, and a short write reaches the
	// client as a truncated body. Nothing here can correct that, because the
	// status and the headers left already.
	_, _ = writer.Write([]byte(contents))

	return h
}

// Terminate runs what is left after the server sent the response.
func (h *RequestHandler) Terminate(
	request contract.ServerRequestContract,
	response contract.ResponseContract,
) {
	h.responseSentHandler.ResponseSent(request, response)
}

// Run handles the request, sends the response, and terminates.
func (h *RequestHandler) Run(request contract.ServerRequestContract, writer http.ResponseWriter) {
	response := h.sendingResponseHandler.SendingResponse(request, h.Handle(request))

	h.Send(response, writer)
	h.Terminate(request, response)
}

// CreateResponseFromThrowable returns the response for the failure.
func (h *RequestHandler) CreateResponseFromThrowable(throwable error) contract.ResponseContract {
	content := ""

	if h.debug {
		content = fmt.Sprintf("%v", throwable)
	}

	return h.responseFactory.CreateResponse(content, constant.StatusCodeInternalServerError, nil)
}

// dispatch runs the request through the router, and turns a failure into a
// response.
func (h *RequestHandler) dispatch(request contract.ServerRequestContract) (response contract.ResponseContract) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		throwable, isError := recovered.(error)
		if !isError {
			throwable = fmt.Errorf("%v", recovered)
		}

		response = h.throwableCaughtHandler.ThrowableCaught(
			request,
			h.CreateResponseFromThrowable(throwable),
			throwable,
		)
	}()

	result := h.requestReceivedHandler.RequestReceived(request)
	if result.IsResponse() {
		return result.GetResponse()
	}

	return h.router.Dispatch(result.GetRequest())
}
