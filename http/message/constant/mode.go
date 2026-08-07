/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package constant

import "slices"

type Mode string

// The Mode values that the framework knows.
const (
	ModeRead            Mode = "r"
	ModeReadWrite       Mode = "r+"
	ModeWrite           Mode = "w"
	ModeWriteRead       Mode = "w+"
	ModeWriteEnd        Mode = "a"
	ModeWriteReadEnd    Mode = "a+"
	ModeCreateWrite     Mode = "x"
	ModeCreateWriteRead Mode = "x+"
	ModeWriteCreate     Mode = "c"
	ModeWriteReadCreate Mode = "c+"
	ModeCloseOnExec     Mode = "e"
)

// readableModes holds each mode that a reader opens a stream in.
var readableModes = []Mode{
	ModeRead,
	ModeReadWrite,
	ModeWriteRead,
	ModeWriteReadEnd,
	ModeCreateWriteRead,
	ModeWriteReadCreate,
}

// writableModes holds each mode that a writer opens a stream in.
var writableModes = []Mode{
	ModeReadWrite,
	ModeWrite,
	ModeWriteRead,
	ModeWriteEnd,
	ModeWriteReadEnd,
	ModeCreateWrite,
	ModeCreateWriteRead,
	ModeWriteCreate,
	ModeWriteReadCreate,
}

// IsReadable reports whether a reader reads a stream that opened in the mode.
func (m Mode) IsReadable() bool {
	return slices.Contains(readableModes, m)
}

// IsWritable reports whether a writer writes a stream that opened in the mode.
func (m Mode) IsWritable() bool {
	return slices.Contains(writableModes, m)
}
