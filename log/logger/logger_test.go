/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package logger_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/log/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/log/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/log/logger"
	logexception "github.com/valkyrjaio/valkyrja-go/v26/log/throwable/exception"
)

const logMessage = "the message"

func TestTheStreamLoggerWritesEachSeverity(t *testing.T) {
	t.Parallel()

	tests := map[constant.LogLevel]func(contract.LoggerContract, string, map[string]any){
		constant.LogLevelDebug:     contract.LoggerContract.Debug,
		constant.LogLevelInfo:      contract.LoggerContract.Info,
		constant.LogLevelNotice:    contract.LoggerContract.Notice,
		constant.LogLevelWarning:   contract.LoggerContract.Warning,
		constant.LogLevelError:     contract.LoggerContract.Error,
		constant.LogLevelCritical:  contract.LoggerContract.Critical,
		constant.LogLevelAlert:     contract.LoggerContract.Alert,
		constant.LogLevelEmergency: contract.LoggerContract.Emergency,
	}

	for level, write := range tests {
		written := &strings.Builder{}

		write(logger.NewStreamLogger(written), logMessage, nil)

		expected := strings.ToUpper(string(level)) + ": " + logMessage + "\n"
		if written.String() != expected {
			t.Errorf("the logger must write %q, but wrote: %q", expected, written.String())
		}
	}
}

func TestTheStreamLoggerWritesTheContextAsJson(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	logger.NewStreamLogger(written).Info(logMessage, map[string]any{"one": "two"})

	if written.String() != `INFO: `+logMessage+` {"one":"two"}`+"\n" {
		t.Errorf("the logger must write the context as JSON, but wrote: %q", written.String())
	}
}

func TestTheStreamLoggerWritesAnEmptyContextForDataThatNoEncoderRenders(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	// A channel is a value that no JSON encoder renders.
	logger.NewStreamLogger(written).Info(logMessage, map[string]any{"one": make(chan int)})

	if !strings.HasSuffix(written.String(), "{}\n") {
		t.Errorf("the logger must write an empty context, but wrote: %q", written.String())
	}
}

func TestTheStreamLoggerWritesAFailureAtTheErrorSeverity(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	logger.NewStreamLogger(written).Throwable(errors.New("the failure"), logMessage)

	if !strings.HasPrefix(written.String(), "ERROR: "+logMessage) {
		t.Errorf("the logger must write the failure at the error severity, but wrote: %q", written.String())
	}

	if !strings.Contains(written.String(), "the failure") {
		t.Errorf("the logger must write what went wrong, but wrote: %q", written.String())
	}
}

func TestLogWritesTheMessageAtTheSeverityThatTheCallerNames(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	err := logger.NewStreamLogger(written).Log(constant.LogLevelWarning, logMessage, nil)
	if err != nil {
		t.Fatalf("a severity that the framework reports must be written, but reported: %v", err)
	}

	if written.String() != "WARNING: "+logMessage+"\n" {
		t.Errorf("the logger must write at the severity, but wrote: %q", written.String())
	}
}

func TestLogReportsASeverityThatTheFrameworkDoesNotReport(t *testing.T) {
	t.Parallel()

	written := &strings.Builder{}

	err := logger.NewStreamLogger(written).Log("verbose", logMessage, nil)

	if _, isInvalid := errors.AsType[*logexception.LogInvalidLogLevelError](err); !isInvalid {
		t.Errorf("a severity that the framework does not report must be reported, but returned: %v", err)
	}

	if written.String() != "" {
		t.Errorf("a message at an unknown severity must not be written, but wrote: %q", written.String())
	}
}

func TestTheNullLoggerWritesNothing(t *testing.T) {
	t.Parallel()

	built := logger.NewNullLogger()

	built.Throwable(errors.New("the failure"), logMessage)
	built.Debug(logMessage, nil)
	built.Info(logMessage, nil)
	built.Notice(logMessage, nil)
	built.Warning(logMessage, nil)
	built.Error(logMessage, nil)
	built.Critical(logMessage, nil)
	built.Alert(logMessage, nil)
	built.Emergency(logMessage, nil)

	if built.Log(constant.LogLevelInfo, logMessage, nil) != nil {
		t.Error("a severity that the framework reports must be taken, but was not")
	}

	if built.Log("verbose", logMessage, nil) == nil {
		t.Error("a severity that the framework does not report must be reported, but was not")
	}
}
