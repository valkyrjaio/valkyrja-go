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
parameters, the parsed body, the parsed JSON, the uploaded files, and the
attributes that the framework fills.

`NewJsonServerRequest` parses the body as JSON as it builds the request, where
the content type states JSON:

```go
built := request.NewJsonServerRequest(uri, constant.RequestMethodPost, body, headers)

name := built.GetParsedJson().Get("name")
```

The other ports declare a `JsonServerRequest` that extends `ServerRequest`. A
promoted `With` copies only the embedded struct, so an inherited one would return
a plain server request and drop the parsed JSON. One struct carries both shapes
here, and the constructor is what parses.

Warning: a body that no decoder reads leaves the parsed JSON empty. A request has
no return to carry the failure, so whatever reads the body reports what is
missing instead.

`ResponseFactory` builds each kind, and the component publishes it so a handler
resolves one rather than building a response itself.

### A `With` Method Cannot Report a Failure

Every port returns the contract from a `With` method, so there is no slot for an
error. Where an argument can be invalid, the receiver keeps its own value, and a
separate function validates:

```go
// Right — the caller checks first, then sets.
if !constant.IsValidPort(port) {
	return exception.NewHttpUriInvalidPortError(port)
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

## Structures

A structure names the fields that a request carries and that a response returns,
so a route states its shape once rather than in each handler.

The other ports spell this segment `Struct`. `struct` is a Go keyword, so it
cannot be a package name, and this port spells the word in full.

```go
requestStructure := structure.NewJsonRequestStructure("users.create", "name", "email")

data := requestStructure.GetDataFromRequest(request)

if requestStructure.DetermineIfRequestContainsExtraData(request) {
	// The client sent a field that the route does not name.
}
```

| Constructor                     | Reads                     |
| :------------------------------ | :------------------------ |
| `NewQueryRequestStructure`      | The query parameters      |
| `NewParsedBodyRequestStructure` | The parsed body           |
| `NewJsonRequestStructure`       | The parsed JSON body      |
| `NewRequestStructure`           | Whatever the source names |

The other ports declare an abstract `RequestStruct` and override the two methods
that read a request. Go has no abstract type and no method override, so a
structure holds the function that names the collection to read.

A response structure maps the name a field carries inside the application to the
name that a client reads, so a rename inside never reaches a client:

```go
responseStructure := structure.NewResponseStructure("users.show", map[string]string{
	"name": "full_name",
})

shaped := responseStructure.GetStructuredData(data, true)
```

A field that the data does not carry reads as nil where `includeAll` is true, and
it is left out where `includeAll` is false. A field that the structure does not
name never reaches a client.

## The Client

The client sends a request to another server and returns its response.

| Type         | Reaches                                        |
| :----------- | :--------------------------------------------- |
| `Client`     | The server, over Go's own HTTP client          |
| `NullClient` | No server, and returns an empty response       |
| `LogClient`  | No server, and records what it would have sent |

```go
response, err := client.NewClient(nil).SendRequest(request)
if err != nil {
	return err
}
```

`NewClient` takes Go's own `*http.Client`, so a caller states its own timeout and
its own transport. A caller that passes nil gets the default one, which states no
timeout.

The client reports a failure where it cannot reach the server, and where the
request or the answer is not one that a message carries. A header that the
framework does not accept is dropped rather than fatal, because a server that
answered with one has still answered.

An application that must not reach the network sends through `NullClient`, and a
developer who wants to read what an application would have sent uses `LogClient`.

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
| `valkyrja.http.message.factory.ResponseFactoryContract`        | The response factory        |
| `valkyrja.http.message.response.ResponseContract`              | The response of the request |
| `valkyrja.http.routing.processor.ProcessorContract`            | The route processor         |
| `valkyrja.http.routing.collection.RouteCollectionContract`     | Every route                 |
| `valkyrja.http.routing.matcher.MatcherContract`                | The matcher                 |
| `valkyrja.http.routing.url.UrlContract`                        | The URL generator           |
| `valkyrja.http.routing.factory.RoutingResponseFactoryContract` | The route redirect factory  |
| `valkyrja.http.routing.dispatcher.RouterContract`              | The router                  |
| `valkyrja.http.routing.data.RouteContract`                     | The route that matched      |
| `valkyrja.http.server.handler.RequestHandlerContract`          | The entry point             |

The processor reads each route before the collection files it, because the
matcher reads the regular expression that the processor builds.

The other ports read a route from an annotation on a controller. Go has no
annotation, so a route provider returns its routes as a literal slice that
`sindri` reads from the source.
