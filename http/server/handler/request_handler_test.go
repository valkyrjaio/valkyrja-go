/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package handler_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	containercontract "github.com/valkyrjaio/valkyrja-go/v26/container/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/container/manager"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/request"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
	middlewaredata "github.com/valkyrjaio/valkyrja-go/v26/http/middleware/data"
	middlewarehandler "github.com/valkyrjaio/valkyrja-go/v26/http/middleware/handler"
	serverconstant "github.com/valkyrjaio/valkyrja-go/v26/http/server/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/server/handler"
)

const bodyText = "the body"

// responseFactoryFixture builds the responses that the handler asks for.
type responseFactoryFixture struct {
	contract.ResponseFactoryContract
}

// CreateResponse builds a response that carries the content.
func (f *responseFactoryFixture) CreateResponse(
	content string,
	statusCode constant.StatusCode,
	headers contract.HeaderCollectionContract,
) contract.ResponseContract {
	return response.NewResponseFromContent(content, statusCode, headers)
}

// routerFixture returns the response it holds, or panics with what it holds.
type routerFixture struct {
	response  contract.ResponseContract
	panicWith any
}

// Dispatch returns the response, or panics.
func (r *routerFixture) Dispatch(_ contract.ServerRequestContract) contract.ResponseContract {
	if r.panicWith != nil {
		panic(r.panicWith)
	}

	return r.response
}

// DispatchRoute returns what Dispatch returns.
func (r *routerFixture) DispatchRoute(
	request contract.ServerRequestContract,
	_ contract.RouteContract,
) contract.ResponseContract {
	return r.Dispatch(request)
}

// newHandler builds a request handler over the router, in the debug mode.
func newHandler(
	container containercontract.ContainerContract,
	router contract.RouterContract,
	debug bool,
) *handler.RequestHandler {
	return handler.NewRequestHandler(
		container,
		router,
		middlewarehandler.NewRequestReceivedHandler(container),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewSendingResponseHandler(container),
		middlewarehandler.NewResponseSentHandler(container),
		&responseFactoryFixture{},
		debug,
	)
}

func TestHandleReturnsWhatTheRouterReturns(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	built := response.NewResponseFromContent(bodyText, constant.StatusCodeOk, nil)

	got := newHandler(container, &routerFixture{response: built}, false).
		Handle(request.NewServerRequest(nil, "", nil, nil))

	if got.GetBody().String() != bodyText {
		t.Errorf("Handle must return the response of the router, but is: %q", got.GetBody().String())
	}
}

func TestHandleBindsTheResponse(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	built := response.NewResponseFromContent(bodyText, constant.StatusCodeOk, nil)

	newHandler(container, &routerFixture{response: built}, false).
		Handle(request.NewServerRequest(nil, "", nil, nil))

	bound, err := container.GetSingleton(serverconstant.ResponseContractServiceID)
	if err != nil {
		t.Fatalf("the handler must bind the response, but reported: %v", err)
	}

	if _, isResponse := bound.(contract.ResponseContract); !isResponse {
		t.Errorf("the binding must be a response, but is: %T", bound)
	}
}

func TestHandleTurnsAPanicIntoAServerError(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	got := newHandler(container, &routerFixture{panicWith: errors.New("the handler failed")}, false).
		Handle(request.NewServerRequest(nil, "", nil, nil))

	if got.GetStatusCode() != constant.StatusCodeInternalServerError {
		t.Errorf("a panic must reach the client as 500, but reported: %d", got.GetStatusCode())
	}

	if got.GetBody().String() != "" {
		t.Errorf("a response must state nothing outside debug mode, but is: %q", got.GetBody().String())
	}
}

func TestHandleStatesTheFailureInDebugMode(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	got := newHandler(container, &routerFixture{panicWith: errors.New("the handler failed")}, true).
		Handle(request.NewServerRequest(nil, "", nil, nil))

	if !strings.Contains(got.GetBody().String(), "the handler failed") {
		t.Errorf("a response must state the failure in debug mode, but is: %q", got.GetBody().String())
	}
}

func TestHandleTurnsAPanicOfAnyValueIntoAServerError(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	got := newHandler(container, &routerFixture{panicWith: "a string"}, true).
		Handle(request.NewServerRequest(nil, "", nil, nil))

	if got.GetStatusCode() != constant.StatusCodeInternalServerError {
		t.Errorf("a panic of any value must reach the client as 500, but reported: %d", got.GetStatusCode())
	}

	if !strings.Contains(got.GetBody().String(), "a string") {
		t.Errorf("the response must state what the panic carried, but is: %q", got.GetBody().String())
	}
}

func TestSendWritesTheStatusTheHeadersAndTheBody(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	contentType, err := header.NewHeader(
		constant.HeaderNameContentType,
		value.NewValueFromValue("text/plain"),
	)
	if err != nil {
		t.Fatalf("NewHeader must build the header, but reported: %v", err)
	}

	built := response.NewResponseFromContent(
		bodyText,
		constant.StatusCodeCreated,
		header.NewHeaderCollection(contentType),
	)

	recorder := httptest.NewRecorder()

	newHandler(container, &routerFixture{}, false).Send(built, recorder)

	if recorder.Code != int(constant.StatusCodeCreated) {
		t.Errorf("Send must write the status code, but wrote: %d", recorder.Code)
	}

	if recorder.Header().Get(constant.HeaderNameContentType) != "text/plain" {
		t.Errorf("Send must write the headers, but wrote: %q",
			recorder.Header().Get(constant.HeaderNameContentType))
	}

	if recorder.Body.String() != bodyText {
		t.Errorf("Send must write the body, but wrote: %q", recorder.Body.String())
	}
}

func TestSendWritesNothingForABodyThatNoReaderReads(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	recorder := httptest.NewRecorder()

	newHandler(container, &routerFixture{}, false).
		Send(response.NewEmptyResponse(nil), recorder)

	if recorder.Body.String() != "" {
		t.Errorf("Send must write no body, but wrote: %q", recorder.Body.String())
	}

	if recorder.Code != int(constant.StatusCodeNoContent) {
		t.Errorf("Send must write the status code, but wrote: %d", recorder.Code)
	}
}

func TestRunHandlesSendsAndTerminates(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	built := response.NewResponseFromContent(bodyText, constant.StatusCodeOk, nil)

	recorder := httptest.NewRecorder()

	newHandler(container, &routerFixture{response: built}, false).
		Run(request.NewServerRequest(nil, "", nil, nil), recorder)

	if recorder.Body.String() != bodyText {
		t.Errorf("Run must send the response, but wrote: %q", recorder.Body.String())
	}
}

func TestTerminateRunsTheResponseSentMiddleware(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	// The handler holds no middleware, so this asserts that the call reaches the
	// stage handler and returns rather than panicking.
	newHandler(container, &routerFixture{}, false).Terminate(
		request.NewServerRequest(nil, "", nil, nil),
		response.NewResponse(nil, 0, nil),
	)
}

func TestTheHandlerSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	built := response.NewResponseFromContent(bodyText, constant.StatusCodeOk, nil)

	var handled contract.ExceptionResponseRequestHandlerContract = newHandler(container, &routerFixture{response: built}, false)

	if handled.Handle(request.NewServerRequest(nil, "", nil, nil)).GetBody().String() != bodyText {
		t.Error("the contract must handle the request, but did not")
	}

	if handled.CreateResponseFromThrowable(errors.New("failed")).GetStatusCode() !=
		constant.StatusCodeInternalServerError {
		t.Error("the contract must turn a failure into a server error, but did not")
	}
}

func TestSendWritesNoBodyWhereNoReaderReadsIt(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)

	// A body that opened in a write-only mode reports a failure on the read, so
	// the status and the headers still reach the client and the body does not.
	built := response.NewResponse(
		stream.NewStream(bodyText, constant.ModeWrite),
		constant.StatusCodeOk,
		nil,
	)

	recorder := httptest.NewRecorder()

	newHandler(container, &routerFixture{}, false).Send(built, recorder)

	if recorder.Code != int(constant.StatusCodeOk) {
		t.Errorf("Send must write the status code, but wrote: %d", recorder.Code)
	}

	if recorder.Body.String() != "" {
		t.Errorf("Send must write no body, but wrote: %q", recorder.Body.String())
	}
}

// endingRequestReceivedFixture is a request-received middleware that ends the
// run with the response it holds.
type endingRequestReceivedFixture struct {
	response contract.ResponseContract
}

// RequestReceived ends the run with the response.
func (m *endingRequestReceivedFixture) RequestReceived(
	_ contract.ServerRequestContract,
	_ contract.RequestReceivedHandlerContract,
) contract.RequestReceivedResultContract {
	return middlewaredata.NewRequestReceivedResponse(m.response)
}

func TestHandleEndsWhereTheRequestReceivedMiddlewareReturnsAResponse(t *testing.T) {
	t.Parallel()

	container := manager.NewContainer(nil)
	ending := response.NewResponseFromContent("ended", constant.StatusCodeAccepted, nil)

	container.Bind("Valkyrja.Tests.Http.EndingMiddleware",
		func(_ containercontract.ContainerContract, _ []any) any {
			return &endingRequestReceivedFixture{response: ending}
		})

	built := handler.NewRequestHandler(
		container,
		&routerFixture{panicWith: errors.New("the router must not run")},
		middlewarehandler.NewRequestReceivedHandler(container, "Valkyrja.Tests.Http.EndingMiddleware"),
		middlewarehandler.NewThrowableCaughtHandler(container),
		middlewarehandler.NewSendingResponseHandler(container),
		middlewarehandler.NewResponseSentHandler(container),
		&responseFactoryFixture{},
		false,
	)

	got := built.Handle(request.NewServerRequest(nil, "", nil, nil))

	if got.GetBody().String() != "ended" {
		t.Errorf("the middleware must end the run before the router, but the response is: %q",
			got.GetBody().String())
	}
}

func TestTheHandlerSatisfiesTheExceptionResponseContract(t *testing.T) {
	t.Parallel()

	var built contract.ExceptionResponseRequestHandlerContract = newHandler(manager.NewContainer(nil), nil, false)

	response := built.CreateResponseFromThrowable(errors.New("the handler failed"))

	if response.GetStatusCode() != constant.StatusCodeInternalServerError {
		t.Errorf("a failure must reach the client as a 500, but reached it as: %d", response.GetStatusCode())
	}
}
