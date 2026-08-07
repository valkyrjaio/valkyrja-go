/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package message_test

import (
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/header/value"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/stream"
)

const acceptName = "Accept"

// newAcceptHeader builds an `Accept` header and fails the test where the header
// cannot be built.
func newAcceptHeader(t *testing.T, values ...contract.ValueContract) contract.HeaderContract {
	t.Helper()

	built, err := header.NewHeader(acceptName, values...)
	if err != nil {
		t.Fatalf("NewHeader must build the header, but reported: %v", err)
	}

	return built
}

func TestNewMessageTakesEachDefault(t *testing.T) {
	t.Parallel()

	built := message.NewMessage("", nil, nil)

	if built.GetProtocolVersion() != constant.ProtocolVersionV11 {
		t.Errorf("the protocol version must default to 1.1, but is: %q", built.GetProtocolVersion())
	}

	if len(built.GetHeaders().GetAll()) != 0 {
		t.Error("the headers must default to empty, but did not")
	}

	if built.GetBody().GetSize() != 0 {
		t.Error("the body must default to empty, but did not")
	}
}

func TestNewMessageHoldsWhatItReceives(t *testing.T) {
	t.Parallel()

	headers := header.NewHeaderCollection(newAcceptHeader(t))
	body := stream.NewStream("the body", constant.ModeReadWrite)

	built := message.NewMessage(constant.ProtocolVersionV2, headers, body)

	if built.GetProtocolVersion() != constant.ProtocolVersionV2 {
		t.Errorf("the protocol version must be the one given, but is: %q", built.GetProtocolVersion())
	}

	if !built.GetHeaders().Has(acceptName) {
		t.Error("the headers must be the ones given, but were not")
	}

	if built.GetBody().GetSize() != len("the body") {
		t.Error("the body must be the one given, but was not")
	}
}

func TestEachSetterReplacesWhatTheMessageHolds(t *testing.T) {
	t.Parallel()

	built := message.NewMessage("", nil, nil)

	headers := header.NewHeaderCollection(newAcceptHeader(t))
	body := stream.NewStream("the body", constant.ModeReadWrite)

	built.SetProtocolVersion(constant.ProtocolVersionV3)
	built.SetHeaders(headers)
	built.SetBody(body)

	if built.GetProtocolVersion() != constant.ProtocolVersionV3 {
		t.Errorf("SetProtocolVersion must replace the version, but it is: %q", built.GetProtocolVersion())
	}

	if !built.GetHeaders().Has(acceptName) {
		t.Error("SetHeaders must replace the headers, but did not")
	}

	if built.GetBody().GetSize() != len("the body") {
		t.Error("SetBody must replace the body, but did not")
	}
}

func TestSetBodyRewindsTheBody(t *testing.T) {
	t.Parallel()

	body := stream.NewStream("the body", constant.ModeReadWrite)

	_, err := body.Read(3)
	if err != nil {
		t.Fatalf("Read must read the bytes, but reported: %v", err)
	}

	built := message.NewMessage("", nil, nil)
	built.SetBody(body)

	position, err := built.GetBody().Tell()
	if err != nil {
		t.Fatalf("Tell must return the position, but reported: %v", err)
	}

	if position != 0 {
		t.Errorf("SetBody must rewind the body, but the pointer is at: %d", position)
	}
}

func TestInjectHeaderAddsAHeaderThatTheCollectionDoesNotHold(t *testing.T) {
	t.Parallel()

	headers := message.InjectHeader(
		newAcceptHeader(t, value.NewValueFromValue("text/html")),
		header.NewHeaderCollection(),
		false,
	)

	if headers.GetHeaderLine(acceptName) != "text/html" {
		t.Errorf("InjectHeader must add the header, but the line is: %q", headers.GetHeaderLine(acceptName))
	}
}

func TestInjectHeaderMergesIntoAHeaderOfTheSameName(t *testing.T) {
	t.Parallel()

	existing := header.NewHeaderCollection(
		newAcceptHeader(t, value.NewValueFromValue("text/html")),
	)

	headers := message.InjectHeader(
		newAcceptHeader(t, value.NewValueFromValue("application/json")),
		existing,
		false,
	)

	if headers.GetHeaderLine(acceptName) != "text/html, application/json" {
		t.Errorf("InjectHeader must merge the values, but the line is: %q", headers.GetHeaderLine(acceptName))
	}
}

func TestInjectHeaderReplacesWhereOverrideIsTrue(t *testing.T) {
	t.Parallel()

	existing := header.NewHeaderCollection(
		newAcceptHeader(t, value.NewValueFromValue("text/html")),
	)

	headers := message.InjectHeader(
		newAcceptHeader(t, value.NewValueFromValue("application/json")),
		existing,
		true,
	)

	if headers.GetHeaderLine(acceptName) != "application/json" {
		t.Errorf("InjectHeader must replace the header, but the line is: %q", headers.GetHeaderLine(acceptName))
	}
}
