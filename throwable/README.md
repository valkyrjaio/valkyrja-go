# Throwable

## Introduction

The Throwable component holds the root contract that every framework error
satisfies, and the two base errors that every component error derives from.

Go has no exceptions. An error is a value, and a function returns it. This port
therefore spells a concrete error `*Error` rather than `*Exception`, which is
foreign to Go, and no component throws.

## The Root Contract

`throwable/contract.ValkyrjaThrowable` is the root contract:

```go
type ValkyrjaThrowable interface {
	error

	// GetTraceCode returns a stable identifier for the site that raised the
	// error.
	GetTraceCode() string
}
```

Each component declares a contract of its own that embeds this one, and each
concrete error satisfies the contract of its component.

A component contract adds a marker method, because Go compares an interface by
its method set. Without the marker, a component contract and the root contract
are the same type, and a caller cannot ask which component raised an error.

```go
// Right — the marker makes the contract a type of its own.
type HttpThrowable interface {
	throwablecontract.ValkyrjaThrowable

	IsHttpThrowable() bool
}
```

## The Two Base Errors

The other ports declare these two abstract. Go has no abstract type, so each base
is a struct that a component error embeds. A component never returns a base error
on its own.

| Type                           | Holds                                        |
| :----------------------------- | :------------------------------------------- |
| `ValkyrjaRuntimeError`         | A failure that the caller cannot correct     |
| `ValkyrjaInvalidArgumentError` | A failure that the caller's own value caused |

Each one records its message, its cause, and the call site that raised it.
`errors.Is` and `errors.As` read the cause, because every base implements
`Unwrap`.

```go
// Right — the concrete error embeds its component's base, which embeds one of
// these two.
type HttpUriInvalidPortError struct {
	HttpUriInvalidArgumentError

	port int
}
```

## The Trace Code

`GetTraceCode` returns a stable identifier for the site that raised an error. Two
errors raised at the same site share a code, and a reader quotes the code to
locate the site.

The other ports hash the error name and the stack with MD5. This port hashes the
recorded frames with SHA-256, because `gosec` rejects MD5. The code is an opaque
identifier, so it only has to be stable — no port can match another port's value
in any case.

## Service Registration

The component registers no service. It holds contracts and base types alone.
