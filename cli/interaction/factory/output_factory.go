/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package factory builds each kind of output that a command writes to.
package factory

import (
	"io"
	"os"

	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/data"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/interaction/output"
)

// filePermissions is what a new output file is created with. The owner reads and
// writes the file, and every other user reads it.
const filePermissions = 0o644

type OutputFactory struct {
	config contract.CliInteractionConfigContract
	writer io.Writer
}

// NewOutputFactory builds the factory over a configuration.
func NewOutputFactory(config contract.CliInteractionConfigContract) *OutputFactory {
	if config == nil {
		config = data.NewCliInteractionConfig()
	}

	return &OutputFactory{
		config: config,
		writer: os.Stdout,
	}
}

// NewOutputFactoryForWriter builds the factory over a configuration and a
// writer.
func NewOutputFactoryForWriter(config contract.CliInteractionConfigContract, writer io.Writer) *OutputFactory {
	built := NewOutputFactory(config)
	built.writer = writer

	return built
}

// CreateOutput builds an output that writes through the default writers.
func (f *OutputFactory) CreateOutput(
	exitCode constant.ExitCode,
	messages ...contract.MessageContract,
) contract.OutputContract {
	return f.CreateStreamOutput(f.writer, exitCode, messages...)
}

// CreateEmptyOutput builds an output that writes nothing.
func (f *OutputFactory) CreateEmptyOutput(
	exitCode constant.ExitCode,
	messages ...contract.MessageContract,
) contract.EmptyOutputContract {
	return f.build(nil, exitCode, messages)
}

// CreatePlainOutput builds an output that applies no format.
func (f *OutputFactory) CreatePlainOutput(
	exitCode constant.ExitCode,
	messages ...contract.MessageContract,
) contract.PlainOutputContract {
	return f.build([]contract.WriterContract{output.NewPlainWriter(f.writer)}, exitCode, messages)
}

// CreateFileOutput builds an output that writes to the file.
func (f *OutputFactory) CreateFileOutput(
	filepath string,
	exitCode constant.ExitCode,
	messages ...contract.MessageContract,
) contract.FileOutputContract {
	// The caller names the file, the same way it does in every other port. An
	// application that takes the name from a caller validates it before it
	// reaches here.
	//nolint:gosec // The path is the application's own, not a client's.
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermissions)
	if err != nil {
		return f.CreateEmptyOutput(exitCode, messages...)
	}

	return f.CreateStreamOutput(file, exitCode, messages...)
}

// CreateStreamOutput builds an output that writes to the writer.
func (f *OutputFactory) CreateStreamOutput(
	writer io.Writer,
	exitCode constant.ExitCode,
	messages ...contract.MessageContract,
) contract.StreamOutputContract {
	return f.build([]contract.WriterContract{output.NewStreamWriter(writer)}, exitCode, messages)
}

// build makes the output that the writers write through, with the settings of
// the configuration.
func (f *OutputFactory) build(
	writers []contract.WriterContract,
	exitCode constant.ExitCode,
	messages []contract.MessageContract,
) contract.OutputContract {
	return output.NewOutput(writers, messages...).
		WithIsInteractive(f.config.IsInteractive()).
		WithIsQuiet(f.config.IsQuiet()).
		WithIsSilent(f.config.IsSilent()).
		WithExitCode(exitCode)
}
