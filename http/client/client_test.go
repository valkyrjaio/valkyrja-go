/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package client_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/client"
	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/request"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/uri"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
	"github.com/valkyrjaio/valkyrja-go/v26/log/logger"
)

const responseBody = "the answer"

// newRequest builds a request to the target.
func newRequest(t *testing.T, target string, body string) contract.RequestContract {
	t.Helper()

	built, err := uri.NewUriFromString(target)
	if err != nil {
		t.Fatalf("the parser must read %q, but reported: %v", target, err)
	}

	headers, err := header.NewHeader(constant.HeaderNameContentType, value.NewValueFromValue("text/plain"))
	if err != nil {
		t.Fatalf("the header must be valid, but reported: %v", err)
	}

	return request.NewRequest(
		built,
		constant.RequestMethodPost,
		stream.NewStream(body, constant.ModeReadWrite),
		header.NewHeaderCollection().WithHeader(headers),
	)
}

func TestTheNullClientReachesNoServer(t *testing.T) {
	t.Parallel()

	var built contract.ClientContract = client.NewNullClient()

	response, err := built.SendRequest(newRequest(t, "https://example.com/", ""))
	if err != nil {
		t.Fatalf("the null client must report no failure, but reported: %v", err)
	}

	if response.GetStatusCode() != constant.StatusCodeNoContent {
		t.Errorf("the null client must return an empty response, but returned: %d", response.GetStatusCode())
	}
}

func TestTheLogClientRecordsTheRequest(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	var built contract.ClientContract = client.NewLogClient(logger.NewStreamLogger(written))

	response, err := built.SendRequest(newRequest(t, "https://example.com/users", ""))
	if err != nil {
		t.Fatalf("the log client must report no failure, but reported: %v", err)
	}

	if response.GetStatusCode() != constant.StatusCodeNoContent {
		t.Error("the log client must return an empty response, but did not")
	}

	for _, part := range []string{"POST", "https://example.com/users"} {
		if !strings.Contains(written.String(), part) {
			t.Errorf("the log client must record %q, but recorded: %q", part, written.String())
		}
	}
}

func TestTheClientSendsTheRequestAndReadsTheAnswer(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, read *http.Request) {
		body, _ := io.ReadAll(read.Body)

		writer.Header().Set("X-Echo-Method", read.Method)
		writer.Header().Set("X-Echo-Body", string(body))
		writer.Header().Set("X-Echo-Content-Type", read.Header.Get(constant.HeaderNameContentType))
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(responseBody))
	}))
	defer server.Close()

	var built contract.ClientContract = client.NewClient(server.Client())

	response, err := built.SendRequest(newRequest(t, server.URL+"/users", "the body"))
	if err != nil {
		t.Fatalf("the client must reach the server, but reported: %v", err)
	}

	if response.GetStatusCode() != constant.StatusCodeCreated {
		t.Errorf("the response must carry the status that the server answered, but carried: %d",
			response.GetStatusCode())
	}

	headers := response.GetHeaders()

	if headers.GetHeaderLine("X-Echo-Method") != "POST" {
		t.Error("the client must send the method of the request, but did not")
	}

	if headers.GetHeaderLine("X-Echo-Body") != "the body" {
		t.Error("the client must send the body of the request, but did not")
	}

	if headers.GetHeaderLine("X-Echo-Content-Type") != "text/plain" {
		t.Error("the client must send each header of the request, but did not")
	}
}

func TestTheClientReportsAServerThatItCannotReach(t *testing.T) {
	t.Parallel()

	built := client.NewClient(&http.Client{Transport: &failingTransport{}})

	_, err := built.SendRequest(newRequest(t, "https://example.com/users", ""))

	failure, isFailure := errors.AsType[*exception.HttpClientRequestFailedError](err)
	if !isFailure {
		t.Fatalf("the client must report a failure of its own, but reported: %v", err)
	}

	if failure.GetUri() != "https://example.com/users" {
		t.Errorf("the failure must name the URI, but named: %q", failure.GetUri())
	}
}

func TestTheClientTakesGoOwnClientWhereTheCallerNamesNone(t *testing.T) {
	t.Parallel()

	if client.NewClient(nil) == nil {
		t.Error("the client must build over Go's own client, but built none")
	}
}

// failingTransport is a transport that reaches no server.
type failingTransport struct{}

// RoundTrip reports the failure.
func (t *failingTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, errors.New("the server is unreachable")
}

// craftedTransport answers with the response that the test built, rather than
// reaching a server.
type craftedTransport struct {
	response *http.Response
}

// RoundTrip returns the crafted response.
func (t *craftedTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return t.response, nil
}

// failingBody is a body that reports a failure rather than its contents.
type failingBody struct{}

// Read reports the failure.
func (b *failingBody) Read(_ []byte) (int, error) {
	return 0, errors.New("the body is unreadable")
}

// Close reports that the body closed.
func (b *failingBody) Close() error {
	return nil
}

func TestTheClientReportsABodyThatItCannotRead(t *testing.T) {
	t.Parallel()

	target, err := url.Parse("https://example.com/users")
	if err != nil {
		t.Fatalf("the parser must read the target, but reported: %v", err)
	}

	built := client.NewClient(&http.Client{Transport: &craftedTransport{response: &http.Response{
		StatusCode: http.StatusOK,
		Body:       &failingBody{},
		Header:     http.Header{},
		Request:    &http.Request{URL: target},
	}}})

	_, err = built.SendRequest(newRequest(t, "https://example.com/users", ""))

	if _, isFailure := errors.AsType[*exception.HttpClientRequestFailedError](err); !isFailure {
		t.Errorf("the client must report a failure of its own, but reported: %v", err)
	}
}

func TestTheClientDropsAHeaderThatTheFrameworkDoesNotAccept(t *testing.T) {
	t.Parallel()

	// A server that answers with a name no header carries has still answered,
	// so the client drops the header rather than the response.
	answered := http.Header{}
	answered["Bad Name"] = []string{"one"}
	answered["X-Good-Name"] = []string{"two"}

	built := client.NewClient(&http.Client{Transport: &craftedTransport{response: &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
		Header:     answered,
	}}})

	response, err := built.SendRequest(newRequest(t, "https://example.com/users", ""))
	if err != nil {
		t.Fatalf("the client must return the response, but reported: %v", err)
	}

	if response.GetHeaders().Has("Bad Name") {
		t.Error("a header that the framework does not accept must be dropped, but was kept")
	}

	if response.GetHeaders().GetHeaderLine("X-Good-Name") != "two" {
		t.Error("every other header must be kept, but was not")
	}
}

func TestTheClientReportsARequestThatItCannotBuild(t *testing.T) {
	t.Parallel()

	tests := map[string]contract.RequestContract{
		"a body that no reader reads": request.NewRequest(
			mustUri(t, "https://example.com/users"),
			constant.RequestMethodPost,
			stream.NewStream("", constant.ModeWrite),
			nil,
		),
		"a method that no request carries": request.NewRequest(
			mustUri(t, "https://example.com/users"),
			"BAD METHOD",
			stream.NewStream("", constant.ModeReadWrite),
			nil,
		),
	}

	for name, sent := range tests {
		_, err := client.NewClient(nil).SendRequest(sent)
		if err == nil {
			t.Errorf("%s must report a failure, but reported none", name)
		}
	}
}

// mustUri reads the target and fails the test where it is not a URI.
func mustUri(t *testing.T, target string) contract.UriContract {
	t.Helper()

	built, err := uri.NewUriFromString(target)
	if err != nil {
		t.Fatalf("the parser must read %q, but reported: %v", target, err)
	}

	return built
}
