/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package data holds the CLI interaction component's configuration.
package data

type CliInteractionConfig struct {
	quiet       bool
	interactive bool
	silent      bool
}

// NewCliInteractionConfig builds the configuration that the framework uses where
// the application states none.
func NewCliInteractionConfig() *CliInteractionConfig {
	return &CliInteractionConfig{interactive: true}
}

// NewCliInteractionConfigFromValues builds the configuration from the values
// that the application states.
func NewCliInteractionConfigFromValues(quiet bool, interactive bool, silent bool) *CliInteractionConfig {
	return &CliInteractionConfig{
		quiet:       quiet,
		interactive: interactive,
		silent:      silent,
	}
}

// IsQuiet reports whether an output writes less.
func (c *CliInteractionConfig) IsQuiet() bool {
	return c.quiet
}

// IsInteractive reports whether an output asks the caller a question.
func (c *CliInteractionConfig) IsInteractive() bool {
	return c.interactive
}

// IsSilent reports whether an output writes nothing.
func (c *CliInteractionConfig) IsSilent() bool {
	return c.silent
}
