/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package message

import (
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
)

// lowestPercentage is how far a command has come before it starts.
const lowestPercentage = 0

// highestPercentage is how far a command has come once it finished.
const highestPercentage = 100

// Progress is a message that reports how far a command has come.
type Progress struct {
	Message

	isComplete bool
	percentage int
}

// NewProgress builds a progress that carries the text, and that reports that the
// command has not started.
func NewProgress(text string) *Progress {
	return &Progress{Message: *NewMessage(text)}
}

// WithText returns a copy of the progress with another text.
//
// The method is declared here rather than promoted from the embedded message,
// because a promoted `With` copies only the embedded struct and would return a
// plain message, dropping every field of the progress.
func (p *Progress) WithText(text string) contract.MessageContract {
	copied := *p
	copied.text = text

	return &copied
}

// WithFormatter returns a copy of the progress with another formatter.
func (p *Progress) WithFormatter(formatter contract.FormatterContract) contract.MessageContract {
	copied := *p
	copied.formatter = formatter

	return &copied
}

// WithoutFormatter returns a copy of the progress with no formatter.
func (p *Progress) WithoutFormatter() contract.MessageContract {
	copied := *p
	copied.formatter = nil

	return &copied
}

// IsComplete reports whether the command finished.
func (p *Progress) IsComplete() bool {
	return p.isComplete
}

// WithIsComplete returns a copy of the progress with another complete flag.
//
// A progress that reports the command finished reports 100 as well, so the two
// never disagree.
func (p *Progress) WithIsComplete(isComplete bool) contract.ProgressContract {
	copied := *p
	copied.isComplete = isComplete

	if isComplete {
		copied.percentage = highestPercentage
	}

	return &copied
}

// GetPercentage returns how far the command has come, from 0 to 100.
func (p *Progress) GetPercentage() int {
	return p.percentage
}

// WithPercentage returns a copy of the progress at another percentage.
//
// A percentage outside the range takes the nearest end of it, because a progress
// that reports more than everything reports nothing a reader can use. A progress
// at 100 reports that the command finished.
func (p *Progress) WithPercentage(percentage int) contract.ProgressContract {
	percentage = min(max(percentage, lowestPercentage), highestPercentage)

	copied := *p
	copied.percentage = percentage
	copied.isComplete = percentage == highestPercentage

	return &copied
}

// A progress satisfies its contract, which the compiler checks at build time
// rather than at run time.
var _ contract.ProgressContract = (*Progress)(nil)
