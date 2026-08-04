# Http

## Introduction

The Http component serves one request. It has three sub-components: message holds
the request, the response, and every part of each, routing matches a request to a
route and runs it, and the server is the entry point.

## Message

### Uri

`NewUri` builds a URI from each of its parts, and `NewUriFromString` reads one
from a string with `net/url`. Each one reports a failure where the port, the
path, or the query string is not one that a URI carries.

```go
built, err := uri.NewUriFromString("https://example.com/users/1?page=2")
if err != nil {
	return err
}
```

A port of zero takes the standard port of the scheme, and the authority leaves a
standard port out.

### Headers

A header holds a name and one or more values, and a value holds one or more
components. `NewHeader` reports a failure where the name is not one that a header
carries.

A header collection reads a header by name without regard to case, and
`GetHeaderLine` returns every value as one comma-separated string.

Warning: a folded header carries `\r\n` followed by whitespace. A bare line feed
is header injection, and the factory rejects it.

### Streams

A stream is the body of a message. It opens in a mode, and a read on a stream
that no reader reads reports a failure, as does a write on one that no writer
writes.

### Requests and Responses

| Constructor                | Builds                                     |
| :------------------------- | :----------------------------------------- |
| `NewRequest`               | A request that a client sends              |
| `NewServerRequest`         | A request that a server received           |
| `NewResponse`              | A response over a body and a status code   |
| `NewTextResponse`          | A response that carries plain text         |
| `NewHtmlResponse`          | A response that carries HTML               |
| `NewEmptyResponse`         | A response that carries nothing            |
| `NewJsonResponseFromData`  | A response that carries data as JSON       |
| `NewJsonpResponseFromData` | A response that wraps the JSON in a call   |
| `NewRedirectResponseToUri` | A response that sends the client elsewhere |

A server request also holds the server parameters, the cookies, the query
parameters, the parsed body, the uploaded files, and the attributes that the
framework fills.

`ResponseFactory` builds each kind, and the component publishes it so a handler
resolves one rather than building a response itself.

### A `With` Method Cannot Report a Failure

Every port returns the contract from a `With` method, so there is no slot for an
error. Where an argument can be invalid, the receiver keeps its own value, and a
separate function validates:

```go
// Right — the caller validates first, then sets.
err := uri.ValidatePort(port)
if err != nil {
	return err
}

built = built.WithPort(port)
```

This holds for `Header.WithName`, `Uri.WithPort`, `Uri.WithPath`,
`Uri.WithQuery`, and `Request.WithRequestTarget`.

## Routing

### A Route

```go
route := data.NewRoute("/users/{id}", "users.show", handler, constant.RequestMethodGet).
	WithParameters(data.NewParameter("id", constant.RegexNum))
```

A handler receives the container and the route that matched, and returns a
response:

```go
type HttpHandlerFunc func(
	container containercontract.ContainerContract,
	route RouteContract,
) ResponseContract
```

A route that names no request method matches `GET`.

### Static and Dynamic Routes

The other ports declare a `DynamicRoute` that extends `Route`. A `With` method
promoted from an embedded struct copies only that struct, so an inherited `With`
on a dynamic route would return a static one and drop its regular expression. One
struct carries both shapes here, and `IsDynamic` reports which it is.

### Regular Expressions

The processor builds the regular expression of a dynamic route from its path.
Go's regular expressions are RE2, which constrains the shape three ways:

- A named group is `(?P<name>…)`, the portable spelling.
- A pattern carries no delimiter, and the `^` and `$` anchors are load-bearing.
  `MatchString` searches, where Java's `matches` implies a whole match on its own.
- RE2 rejects lookahead, lookbehind, and a backreference outright.

Warning: an optional parameter carries its own separator, so its path leaves the
separator out. Write `/users{id?}`, never `/users/{id?}` — the second one needs
the separator, so `/users` does not match it.

### Casting a Parameter

A parameter names the function that converts the value a path filled into it, so
a handler reads the type that the route declared rather than the raw text:

```go
parameter := data.NewParameter("id", constant.RegexNum).
	WithCast(func(value string) (any, error) {
		return strconv.Atoi(value)
	})
```

The other ports carry a `Cast` from the Type component, which this port does not
have. The CLI component spells its own cast the same way.

Warning: the regular expression states the shape of a value, and the cast
converts what the regular expression accepted already. A cast that reports a
failure means the two disagree, which is the developer's error. The parameter
then carries the text as the path held it.

### The Router

The router matches the request path against the static routes first, then against
the dynamic ones. The matched route carries the value of each parameter that the
path filled, so the handler reads them from the route.

A path that matches no route of the request method, but matches one of another,
reports the methods that it does match.

`Url` reads the URL of a named route, and `RoutingResponseFactory` builds a
response that sends the client to one.

## Middleware

| Stage             | Runs                                            |
| :---------------- | :---------------------------------------------- |
| `RequestReceived` | Before the router matches a route               |
| `RouteMatched`    | When the router matches one                     |
| `RouteNotMatched` | When the router matches none                    |
| `RouteDispatched` | After the route handler returns                 |
| `ThrowableCaught` | When something in the request reports a failure |
| `SendingResponse` | Before the server writes the response           |
| `ResponseSent`    | After the client received it                    |

The request-received and route-matched stages return a result contract, because
the TypeScript port returns a union there and Go has no union.

Warning: a handler appends each middleware and never dedupes. A middleware that
is added twice runs twice, and the framework does not correct it.

## The Server

`RequestHandler` is the entry point for one request:

```go
handler.Run(request, writer)
```

Go's `http.ResponseWriter` takes the headers before the status code, and the
status code before the body. A write in another order is silently dropped, so
that order is not a style choice.

A handler that panics reaches the client as a 500 rather than ending the process.
The other ports catch a throw here; Go has no throw, so this recovers a panic
instead, which is one of two places in the framework that does.

Warning: a response states what went wrong only in debug mode. A production
application leaves `DebugMode` false, so a failure never tells a client about the
inside of the application.

## Configuration

The component reads `application/contract.HttpConfigContract`, which names the
middleware that each stage runs, by binding key, and carries the debug mode.

## Service Registration

| Binding key                                                    | Holds                       |
| :------------------------------------------------------------- | :-------------------------- |
| `Valkyrja.Http.Message.Factory.ResponseFactoryContract`        | The response factory        |
| `Valkyrja.Http.Message.Response.ResponseContract`              | The response of the request |
| `Valkyrja.Http.Routing.Processor.ProcessorContract`            | The route processor         |
| `Valkyrja.Http.Routing.Collection.RouteCollectionContract`     | Every route                 |
| `Valkyrja.Http.Routing.Matcher.MatcherContract`                | The matcher                 |
| `Valkyrja.Http.Routing.Url.UrlContract`                        | The URL generator           |
| `Valkyrja.Http.Routing.Factory.RoutingResponseFactoryContract` | The route redirect factory  |
| `Valkyrja.Http.Routing.Dispatcher.RouterContract`              | The router                  |
| `Valkyrja.Http.Routing.Data.RouteContract`                     | The route that matched      |
| `Valkyrja.Http.Server.Handler.RequestHandlerContract`          | The entry point             |

The processor reads each route before the collection files it, because the
matcher reads the regular expression that the processor builds.

The other ports read a route from an annotation on a controller. Go has no
annotation, so a route provider returns its routes as a literal slice that
`sindri` reads from the source.
