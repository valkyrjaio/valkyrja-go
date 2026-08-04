/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package exception

// HttpHeaderInvalidArgumentError is the header sub-component's base
// invalid-argument error.
type HttpHeaderInvalidArgumentError struct {
	HttpInvalidArgumentError
}

// HttpHeaderInvalidNameError reports a header name that no header carries.
type HttpHeaderInvalidNameError struct {
	HttpHeaderInvalidArgumentError

	name string
}

// NewHttpHeaderInvalidNameError builds the error for a header name.
func NewHttpHeaderInvalidNameError(name string) *HttpHeaderInvalidNameError {
	return &HttpHeaderInvalidNameError{
		HttpHeaderInvalidArgumentError: HttpHeaderInvalidArgumentError{
			HttpInvalidArgumentError: NewHttpInvalidArgumentError(
				`"`+name+`" is not valid header name`,
				nil,
			),
		},
		name: name,
	}
}

// GetName returns the header name that the error reports.
func (e *HttpHeaderInvalidNameError) GetName() string {
	return e.name
}

// HttpHeaderInvalidValueError reports a header value that no header carries.
type HttpHeaderInvalidValueError struct {
	HttpHeaderInvalidArgumentError

	value string
}

// NewHttpHeaderInvalidValueError builds the error for a header value.
func NewHttpHeaderInvalidValueError(value string) *HttpHeaderInvalidValueError {
	return &HttpHeaderInvalidValueError{
		HttpHeaderInvalidArgumentError: HttpHeaderInvalidArgumentError{
			HttpInvalidArgumentError: NewHttpInvalidArgumentError(
				`"`+value+`" is not valid header value`,
				nil,
			),
		},
		value: value,
	}
}

// GetValue returns the header value that the error reports.
func (e *HttpHeaderInvalidValueError) GetValue() string {
	return e.value
}

// HttpHeaderInvalidHeaderNameError reports a header that a collection does not
// hold.
type HttpHeaderInvalidHeaderNameError struct {
	HttpHeaderInvalidArgumentError

	name string
}

// NewHttpHeaderInvalidHeaderNameError builds the error for a header name.
func NewHttpHeaderInvalidHeaderNameError(name string) *HttpHeaderInvalidHeaderNameError {
	return &HttpHeaderInvalidHeaderNameError{
		HttpHeaderInvalidArgumentError: HttpHeaderInvalidArgumentError{
			HttpInvalidArgumentError: NewHttpInvalidArgumentError(
				"Header "+name+" does not exist",
				nil,
			),
		},
		name: name,
	}
}

// GetName returns the header name that the error reports.
func (e *HttpHeaderInvalidHeaderNameError) GetName() string {
	return e.name
}
