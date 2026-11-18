/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package message

import (
	"runtime"
	"strings"

	applicationconstant "github.com/valkyrjaio/valkyrja-go/v26/application/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/cli/contract"
)

// headerIndent opens each line that the header draws inside its frame.
const headerIndent = "│   "

type Header struct {
	appName    string
	appVersion string

	icon              string
	valkyrjaVersion   string
	valkyrjaBuildDate string
	goVersion         string
	projectRoot       string
	actionDescription string
	commandName       string
}

// NewHeader builds the header for the application and the command.
func NewHeader(appName string, appVersion string, projectRoot string, route contract.RouteContract) *Header {
	built := &Header{
		appName:           appName,
		appVersion:        appVersion,
		icon:              applicationconstant.Icon,
		valkyrjaVersion:   applicationconstant.Version,
		valkyrjaBuildDate: applicationconstant.VersionBuildDateTime,
		goVersion:         runtime.Version(),
		projectRoot:       projectRoot,
	}

	if route != nil {
		built.actionDescription = route.GetDescription()
		built.commandName = route.GetName()
	}

	return built
}

// WithAppName returns a copy of the header for another application name.
func (h *Header) WithAppName(appName string) *Header {
	copied := *h
	copied.appName = appName

	return &copied
}

// WithAppVersion returns a copy of the header for another application version.
func (h *Header) WithAppVersion(appVersion string) *Header {
	copied := *h
	copied.appVersion = appVersion

	return &copied
}

// WithIcon returns a copy of the header that draws another icon.
func (h *Header) WithIcon(icon string) *Header {
	copied := *h
	copied.icon = icon

	return &copied
}

// WithValkyrjaVersion returns a copy of the header for another framework
// version.
func (h *Header) WithValkyrjaVersion(valkyrjaVersion string) *Header {
	copied := *h
	copied.valkyrjaVersion = valkyrjaVersion

	return &copied
}

// WithValkyrjaBuildDate returns a copy of the header for another framework build
// date.
func (h *Header) WithValkyrjaBuildDate(valkyrjaBuildDate string) *Header {
	copied := *h
	copied.valkyrjaBuildDate = valkyrjaBuildDate

	return &copied
}

// WithGoVersion returns a copy of the header for another toolchain version.
func (h *Header) WithGoVersion(goVersion string) *Header {
	copied := *h
	copied.goVersion = goVersion

	return &copied
}

// WithProjectRoot returns a copy of the header for another project root.
func (h *Header) WithProjectRoot(projectRoot string) *Header {
	copied := *h
	copied.projectRoot = projectRoot

	return &copied
}

// WithActionDescription returns a copy of the header for another description.
func (h *Header) WithActionDescription(actionDescription string) *Header {
	copied := *h
	copied.actionDescription = actionDescription

	return &copied
}

// WithCommandName returns a copy of the header for another command name.
func (h *Header) WithCommandName(commandName string) *Header {
	copied := *h
	copied.commandName = commandName

	return &copied
}

// GetText returns the header as the terminal prints it.
func (h *Header) GetText() string {
	lines := []string{
		"╭── " + h.appName + " v" + h.appVersion,
		"│",
	}

	for line := range strings.SplitSeq(h.icon, "\n") {
		lines = append(lines, headerIndent+line)
	}

	return strings.Join(append(lines,
		"│",
		headerIndent+"Built on Valkyrja v"+h.valkyrjaVersion+" (date: "+h.valkyrjaBuildDate+")",
		headerIndent+"Running on "+h.goVersion,
		headerIndent+h.projectRoot,
		"╰── "+h.actionDescription+" · "+h.commandName,
	), "\n")
}

// GetFormattedText returns the header, which carries no format of its own.
func (h *Header) GetFormattedText() string {
	return h.GetText()
}

// WithText returns a message that carries the text.
func (h *Header) WithText(text string) contract.MessageContract {
	return NewMessage(text)
}

// HasFormatter reports that a header carries no formatter.
func (h *Header) HasFormatter() bool {
	return false
}

// GetFormatter returns nil, because a header carries no formatter.
func (h *Header) GetFormatter() contract.FormatterContract {
	return nil
}

// WithFormatter returns a message that carries the header text and the
// formatter.
func (h *Header) WithFormatter(formatter contract.FormatterContract) contract.MessageContract {
	return NewFormattedMessage(h.GetText(), formatter)
}

// WithoutFormatter returns the header, which carries no formatter already.
func (h *Header) WithoutFormatter() contract.MessageContract {
	return h
}
