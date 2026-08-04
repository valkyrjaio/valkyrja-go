/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

// Package constant holds the HTTP routing component's regular expressions and
// its binding keys.
package constant

// uuidPart is one hexadecimal character of a UUID.
const uuidPart = "[0-9A-Fa-f]"

// ulidVlidChars is the alphabet that a ULID and a VLID are written in.
const ulidVlidChars = "0123456789ABCDEFGHJKMNPQRSTVWXYZabcdefghjkmnpqrstvwxyz"

// The regular expression that each kind of route parameter matches.
//
// Warning: Go's regular expressions are RE2, which holds no lookahead, no
// lookbehind, and no backreference. Every pattern here uses a character class, a
// quantifier, or an alternation only, so each one compiles on every port.
const (
	// RegexPath is the separator between two segments of a path.
	RegexPath = `\/`

	// RegexAny matches any character.
	RegexAny = `.*`

	// RegexNum matches one or more digits.
	RegexNum = `\d+`

	// RegexId matches an identifier, which is a number.
	RegexId = RegexNum

	// RegexSlug matches a slug.
	RegexSlug = `[a-zA-Z0-9-]+`

	// RegexAlpha matches one or more letters.
	RegexAlpha = `[a-zA-Z]+`

	// RegexAlphaLowercase matches one or more lower-case letters.
	RegexAlphaLowercase = `[a-z]+`

	// RegexAlphaUppercase matches one or more upper-case letters.
	RegexAlphaUppercase = `[A-Z]+`

	// RegexAlphaNum matches one or more letters and digits.
	RegexAlphaNum = `[a-zA-Z0-9]+`

	// RegexAlphaNumUnderscore matches one or more word characters.
	RegexAlphaNumUnderscore = `\w+`

	// RegexUuid matches a UUID of any version.
	RegexUuid = uuidPart + `{8}-` + uuidPart + `{4}-` + uuidPart + `{4}-` +
		uuidPart + `{4}-` + uuidPart + `{12}`

	// RegexUlid matches a ULID.
	RegexUlid = `[0-7][` + ulidVlidChars + `]{25}`
)

// The parts that the processor builds a route's regular expression from.
//
// Warning: a pattern carries no delimiter. PHP writes `/^…$/`, because
// `preg_match` needs the slashes; `regexp.Compile` takes a bare pattern, and the
// slashes would match as literal characters.
//
// The anchors are load-bearing rather than decoration. Go's `MatchString`
// searches, where Java's `matches` implies a whole match on its own, so a route
// without `^` and `$` matches every path that carries its own path inside.
const (
	// RegexStart anchors the pattern to the start of the path.
	RegexStart = `^`

	// RegexEnd anchors the pattern to the end of the path.
	RegexEnd = `$`

	// RegexStartCaptureGroup starts a group that the matcher reads back.
	RegexStartCaptureGroup = `(`

	// RegexStartNonCaptureGroup starts a group that the matcher does not read
	// back.
	RegexStartNonCaptureGroup = `(?:`

	// RegexEndCaptureGroup ends a group.
	RegexEndCaptureGroup = `)`

	// RegexEndOptionalCaptureGroup ends a group that the path may leave out.
	RegexEndOptionalCaptureGroup = `)?`

	// RegexStartOptionalCaptureGroup starts the separator that an optional
	// parameter may leave out with its own segment.
	RegexStartOptionalCaptureGroup = RegexStartNonCaptureGroup + RegexPath + RegexEndOptionalCaptureGroup

	// RegexStartCaptureGroupName starts the name of a capture group.
	//
	// Go spells a named group `(?P<name>…)`. Go 1.22 and later also read the
	// `(?<name>…)` spelling that Java and TypeScript write, but `(?P<name>…)` is
	// the portable form: it compiles on every toolchain, and it is the only
	// spelling that Python reads.
	RegexStartCaptureGroupName = `?P<`

	// RegexEndCaptureGroupName ends the name of a capture group.
	RegexEndCaptureGroupName = `>`
)

// The HTTP routing component's binding keys.
const (
	// RouterContractServiceID is the binding key for the router.
	RouterContractServiceID = "Valkyrja.Http.Routing.Dispatcher.RouterContract"

	// RouteCollectionContractServiceID is the binding key for the route
	// collection.
	RouteCollectionContractServiceID = "Valkyrja.Http.Routing.Collection.RouteCollectionContract"

	// RouteContractServiceID is the binding key for the route that matched. The
	// router binds it before it runs the handler, so a handler reads the route.
	RouteContractServiceID = "Valkyrja.Http.Routing.Data.RouteContract"

	// HttpRoutingDataServiceID is the binding key for the routing data.
	HttpRoutingDataServiceID = "Valkyrja.Http.Routing.Data.HttpRoutingData"

	// MatcherContractServiceID is the binding key for the matcher.
	MatcherContractServiceID = "Valkyrja.Http.Routing.Matcher.MatcherContract"

	// ProcessorContractServiceID is the binding key for the processor.
	ProcessorContractServiceID = "Valkyrja.Http.Routing.Processor.ProcessorContract"

	// UrlContractServiceID is the binding key for the URL generator.
	UrlContractServiceID = "Valkyrja.Http.Routing.Url.UrlContract"
)

// RoutingResponseFactoryContractServiceID is the binding key for what builds a
// response that sends the client to a named route.
const RoutingResponseFactoryContractServiceID = "Valkyrja.Http.Routing.Factory.RoutingResponseFactoryContract"
