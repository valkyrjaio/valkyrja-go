/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package client sends a request to another server and returns its response.
package client

import (
	"context"
	"io"
	"net/http"
	"strings"

	logcontract "github.com/valkyrjaio/valkyrja-go/v26/log/contract"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/response"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

type NullClient struct{}

// NewNullClient builds the client.
func NewNullClient() *NullClient {
	return &NullClient{}
}

// SendRequest returns an empty response, and reaches no server.
func (c *NullClient) SendRequest(_ contract.RequestContract) (contract.ResponseContract, error) {
	return response.NewEmptyResponse(nil), nil
}

type LogClient struct {
	logger logcontract.LoggerContract
}

// NewLogClient builds the client over a logger.
func NewLogClient(logger logcontract.LoggerContract) *LogClient {
	return &LogClient{logger: logger}
}

// SendRequest records the request and returns an empty response.
func (c *LogClient) SendRequest(request contract.RequestContract) (contract.ResponseContract, error) {
	c.logger.Info("LogClient request", map[string]any{
		"method": string(request.GetMethod()),
		"uri":    request.GetUri().String(),
	})

	return response.NewEmptyResponse(nil), nil
}

type Client struct {
	client *http.Client
}

// NewClient builds the client over Go's own HTTP client. A caller that passes
// nil gets the default one, which states no timeout.
func NewClient(built *http.Client) *Client {
	if built == nil {
		built = http.DefaultClient
	}

	return &Client{client: built}
}

// SendRequest sends the request and returns what the server answered.
func (c *Client) SendRequest(request contract.RequestContract) (contract.ResponseContract, error) {
	sent, err := c.newRequest(request)
	if err != nil {
		return nil, err
	}

	received, err := c.client.Do(sent)
	if err != nil {
		return nil, exception.NewHttpClientRequestFailedError(request.GetUri().String(), err)
	}

	defer func() {
		_ = received.Body.Close()
	}()

	return newResponse(received)
}

// newRequest builds the request that Go's own client sends.
func (c *Client) newRequest(request contract.RequestContract) (*http.Request, error) {
	body := request.GetBody()

	_ = body.Rewind()

	contents, err := body.GetContents()
	if err != nil {
		return nil, err
	}

	sent, err := http.NewRequestWithContext(
		context.Background(),
		string(request.GetMethod()),
		request.GetUri().String(),
		strings.NewReader(contents),
	)
	if err != nil {
		return nil, exception.NewHttpClientRequestFailedError(request.GetUri().String(), err)
	}

	for _, held := range request.GetHeaders().GetAll() {
		sent.Header.Set(held.GetName(), held.GetHeaderLine())
	}

	return sent, nil
}

// newResponse builds the response from what the server answered.
func newResponse(received *http.Response) (contract.ResponseContract, error) {
	contents, err := io.ReadAll(received.Body)
	if err != nil {
		return nil, exception.NewHttpClientRequestFailedError(received.Request.URL.String(), err)
	}

	return response.NewResponseFromContent(
		string(contents),
		constant.StatusCode(received.StatusCode),
		newHeaders(received.Header),
	), nil
}

// newHeaders builds the header collection from what the server answered.
func newHeaders(received http.Header) contract.HeaderCollectionContract {
	built := contract.HeaderCollectionContract(header.NewHeaderCollection())

	for name, values := range received {
		components := make([]contract.ValueContract, 0, len(values))

		for _, held := range values {
			components = append(components, value.NewValueFromValue(held))
		}

		made, err := header.NewHeader(name, components...)
		if err != nil {
			continue
		}

		built = built.WithHeader(made)
	}

	return built
}
