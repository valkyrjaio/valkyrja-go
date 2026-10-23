/*
 * This file is part of the Valkyrja Framework package.
 *
 * Copyright (c) 2016-present Melech Mizrachi
 *
 * Released under the MIT License. See LICENSE.md for details.
 */

package uri_test

import (
	"errors"
	"testing"

	"github.com/valkyrjaio/valkyrja-go/v26/http/contract"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/constant"
	"github.com/valkyrjaio/valkyrja-go/v26/http/message/uri"
	"github.com/valkyrjaio/valkyrja-go/v26/http/throwable/exception"
)

const (
	exampleHost  = "example.com"
	basicPath    = "/path"
	hostPort8080 = "example.com:8080"
)

// newUri builds a URI and fails the test where a part is invalid.
func newUri(t *testing.T, scheme constant.Scheme, host string, port int, path string) *uri.Uri {
	t.Helper()

	built, err := uri.NewUri(scheme, "", "", host, port, path, "", "")
	if err != nil {
		t.Fatalf("NewUri must build the URI, but reported: %v", err)
	}

	return built
}

func TestNewUriTakesTheStandardPortOfTheScheme(t *testing.T) {
	t.Parallel()

	tests := map[constant.Scheme]int{
		constant.SchemeHttp:  constant.PortHttp,
		constant.SchemeHttps: constant.PortHttps,
		constant.SchemeEmpty: 0,
	}

	for scheme, port := range tests {
		built := newUri(t, scheme, exampleHost, 0, "")

		if built.HasPort() != (port != 0) {
			t.Errorf("HasPort for %q must be %t, but is %t", scheme, port != 0, built.HasPort())
		}
	}
}

func TestNewUriReportsAPortOutsideTheRange(t *testing.T) {
	t.Parallel()

	_, err := uri.NewUri(constant.SchemeHttp, "", "", exampleHost, 70000, "", "", "")

	target, found := errors.AsType[*exception.HttpUriInvalidPortError](err)
	if !found {
		t.Fatalf("NewUri must report an invalid port, but reported: %v", err)
	}

	if target.GetPort() != 70000 {
		t.Errorf("the error must carry the port, but carries: %d", target.GetPort())
	}
}

func TestNewUriReportsAPathThatCarriesAQueryString(t *testing.T) {
	t.Parallel()

	_, err := uri.NewUri(constant.SchemeHttp, "", "", exampleHost, 0, "/path?a=b", "", "")

	target, found := errors.AsType[*exception.HttpUriInvalidPathError](err)
	if !found {
		t.Fatalf("NewUri must report an invalid path, but reported: %v", err)
	}

	if target.GetPath() != "/path?a=b" {
		t.Errorf("the error must carry the path, but carries: %q", target.GetPath())
	}
}

func TestNewUriReportsAQueryStringThatCarriesAFragment(t *testing.T) {
	t.Parallel()

	_, err := uri.NewUri(constant.SchemeHttp, "", "", exampleHost, 0, "", "a=b#part", "")

	target, found := errors.AsType[*exception.HttpUriInvalidQueryError](err)
	if !found {
		t.Fatalf("NewUri must report an invalid query string, but reported: %v", err)
	}

	if target.GetQuery() != "a=b#part" {
		t.Errorf("the error must carry the query string, but carries: %q", target.GetQuery())
	}
}

func TestTheAuthorityLeavesOutAStandardPort(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		scheme constant.Scheme
		port   int
		want   string
	}{
		"the standard http port": {constant.SchemeHttp, constant.PortHttp, exampleHost},
		"the standard https port": {
			constant.SchemeHttps, constant.PortHttps, exampleHost,
		},
		"another port": {constant.SchemeHttp, 8080, exampleHost + ":8080"},
	}

	for name, test := range tests {
		built := newUri(t, test.scheme, exampleHost, test.port, "")

		if built.GetAuthority() != test.want {
			t.Errorf("the authority with %s must be %q, but is %q", name, test.want, built.GetAuthority())
		}
	}
}

func TestTheAuthorityIsEmptyWithNoHost(t *testing.T) {
	t.Parallel()

	if newUri(t, constant.SchemeHttp, "", 0, "").GetAuthority() != "" {
		t.Error("the authority must be empty with no host, but is not")
	}
}

func TestTheAuthorityCarriesTheUserInfo(t *testing.T) {
	t.Parallel()

	built, err := uri.NewUri(constant.SchemeHttp, "user", "secret", exampleHost, 8080, "", "", "")
	if err != nil {
		t.Fatalf("NewUri must build the URI, but reported: %v", err)
	}

	if built.GetAuthority() != "user:secret@example.com:8080" {
		t.Errorf("the authority must carry the user info, but is: %q", built.GetAuthority())
	}

	if built.GetUserInfo() != "user:secret" {
		t.Errorf("GetUserInfo must join the user name and the password, but is: %q", built.GetUserInfo())
	}
}

func TestGetPortIsZeroForAStandardPort(t *testing.T) {
	t.Parallel()

	if newUri(t, constant.SchemeHttps, exampleHost, constant.PortHttps, "").GetPort() != 0 {
		t.Error("GetPort must be zero for a standard port, but is not")
	}

	if newUri(t, constant.SchemeHttp, exampleHost, 8080, "").GetPort() != 8080 {
		t.Error("GetPort must be the port where it is not the standard one, but is not")
	}
}

func TestGetHostPortAndGetSchemeHostPortJoinTheParts(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeHttps, exampleHost, 8443, "")

	if built.GetHostPort() != "example.com:8443" {
		t.Errorf("GetHostPort must join the host and the port, but is: %q", built.GetHostPort())
	}

	if built.GetSchemeHostPort() != "https://example.com:8443" {
		t.Errorf("GetSchemeHostPort must join every part, but is: %q", built.GetSchemeHostPort())
	}
}

func TestGetHostPortLeavesOutAStandardPort(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeHttps, exampleHost, constant.PortHttps, "")

	if built.GetHostPort() != exampleHost {
		t.Errorf("GetHostPort must leave out a standard port, but is: %q", built.GetHostPort())
	}
}

func TestGetSchemeHostPortIsTheHostPortWithNoScheme(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeEmpty, exampleHost, 8080, "")

	if built.GetSchemeHostPort() != hostPort8080 {
		t.Errorf("GetSchemeHostPort must be the host and port with no scheme, but is: %q",
			built.GetSchemeHostPort())
	}
}

func TestIsSecureReportsTheHttpsScheme(t *testing.T) {
	t.Parallel()

	if !newUri(t, constant.SchemeHttps, exampleHost, 0, "").IsSecure() {
		t.Error("IsSecure must be true for the https scheme, but is false")
	}

	if newUri(t, constant.SchemeHttp, exampleHost, 0, "").IsSecure() {
		t.Error("IsSecure must be false for the http scheme, but is true")
	}
}

func TestTheHostIsLowerCased(t *testing.T) {
	t.Parallel()

	if newUri(t, constant.SchemeHttp, "EXAMPLE.COM", 0, "").GetHost() != exampleHost {
		t.Error("the host must be lower case, but is not")
	}
}

func TestAnIpLiteralHostIsNotEncoded(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeHttp, "[2001:db8::1]", 0, "")

	if built.GetHost() != "[2001:db8::1]" {
		t.Errorf("an IP literal must be left as it is, but is: %q", built.GetHost())
	}
}

func TestThePathCollapsesItsLeadingSeparators(t *testing.T) {
	t.Parallel()

	if newUri(t, constant.SchemeHttp, exampleHost, 0, "///path").GetPath() != basicPath {
		t.Error("the path must collapse its leading separators, but did not")
	}
}

func TestARelativePathKeepsItsShape(t *testing.T) {
	t.Parallel()

	if newUri(t, constant.SchemeHttp, exampleHost, 0, "path/other").GetPath() != "path/other" {
		t.Error("a relative path must keep its shape, but did not")
	}
}

func TestEachPartIsPercentEncoded(t *testing.T) {
	t.Parallel()

	built, err := uri.NewUri(constant.SchemeHttp, "", "", exampleHost, 0, "/a b", "q=a b", "a b")
	if err != nil {
		t.Fatalf("NewUri must build the URI, but reported: %v", err)
	}

	if built.GetPath() != "/a%20b" {
		t.Errorf("the path must be encoded, but is: %q", built.GetPath())
	}

	if built.GetQuery() != "q=a%20b" {
		t.Errorf("the query string must be encoded, but is: %q", built.GetQuery())
	}

	if built.GetFragment() != "a%20b" {
		t.Errorf("the fragment must be encoded, but is: %q", built.GetFragment())
	}
}

func TestAnEncodedTripletKeepsItsMeaning(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeHttp, exampleHost, 0, "/a%2fb")

	if built.GetPath() != "/a%2Fb" {
		t.Errorf("a triplet must keep its meaning and take upper case, but is: %q", built.GetPath())
	}
}

func TestALiteralPercentSignIsEncoded(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"/100%":   "/100%25",
		"/a%zzb":  "/a%25zzb",
		"/a%2":    "/a%252",
		"/a%2Gb":  "/a%252Gb",
		"/50%%20": "/50%25%20",
	}

	for path, want := range tests {
		if newUri(t, constant.SchemeHttp, exampleHost, 0, path).GetPath() != want {
			t.Errorf("the path %q must encode to %q, but is %q",
				path, want, newUri(t, constant.SchemeHttp, exampleHost, 0, path).GetPath())
		}
	}
}

func TestTheQueryStringDropsItsLeadingQuestionMarks(t *testing.T) {
	t.Parallel()

	built, err := uri.NewUri(constant.SchemeHttp, "", "", exampleHost, 0, "", "??a=b", "##part")
	if err != nil {
		t.Fatalf("NewUri must build the URI, but reported: %v", err)
	}

	if built.GetQuery() != "a=b" {
		t.Errorf("the query string must drop its leading question marks, but is: %q", built.GetQuery())
	}

	if built.GetFragment() != "part" {
		t.Errorf("the fragment must drop its leading number signs, but is: %q", built.GetFragment())
	}
}

func TestStringRendersEveryPart(t *testing.T) {
	t.Parallel()

	built, err := uri.NewUri(constant.SchemeHttps, "user", "secret", exampleHost, 8443, "path", "a=b", "part")
	if err != nil {
		t.Fatalf("NewUri must build the URI, but reported: %v", err)
	}

	want := "https://user:secret@example.com:8443/path?a=b#part"

	if built.String() != want {
		t.Errorf("String must be %q, but is %q", want, built.String())
	}
}

func TestStringRendersAUriWithNoAuthority(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeEmpty, "", 0, "/path")

	if built.String() != basicPath {
		t.Errorf("String must be the path alone, but is: %q", built.String())
	}
}

func TestEachWithMethodReturnsACopy(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeHttp, exampleHost, 0, basicPath)

	tests := map[string]struct {
		apply func() contract.UriContract
		want  string
	}{
		"WithScheme": {
			func() contract.UriContract { return built.WithScheme(constant.SchemeHttps) },
			// The receiver already holds port 80 from its own scheme, and
			// WithScheme fills the port only where the URI names none. Port 80
			// is not the standard port of https, so the authority states it.
			"https://example.com:80/path",
		},
		"WithHost":     {func() contract.UriContract { return built.WithHost("other.com") }, "http://other.com/path"},
		"WithPort":     {func() contract.UriContract { return built.WithPort(8080) }, "http://example.com:8080/path"},
		"WithPath":     {func() contract.UriContract { return built.WithPath("/other") }, "http://example.com/other"},
		"WithQuery":    {func() contract.UriContract { return built.WithQuery("a=b") }, "http://example.com/path?a=b"},
		"WithFragment": {func() contract.UriContract { return built.WithFragment("part") }, "http://example.com/path#part"},
		"WithUsername": {func() contract.UriContract { return built.WithUsername("user") }, "http://user@example.com/path"},
	}

	for name, test := range tests {
		if test.apply().String() != test.want {
			t.Errorf("%s must render %q, but rendered %q", name, test.want, test.apply().String())
		}
	}

	if built.String() != "http://example.com/path" {
		t.Error("each With method must leave the receiver unchanged, but did not")
	}
}

func TestWithSchemeTakesTheStandardPortOfTheNewScheme(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeEmpty, exampleHost, 0, "")

	if built.WithScheme(constant.SchemeHttps).String() != "https://example.com" {
		t.Errorf("WithScheme must take the standard port, but rendered: %q",
			built.WithScheme(constant.SchemeHttps).String())
	}
}

func TestWithPasswordAndWithUserInfoHoldBothParts(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeHttp, exampleHost, 0, "")

	withBoth := built.WithUsername("user").WithPassword("secret")

	if withBoth.GetUserInfo() != "user:secret" {
		t.Errorf("WithPassword must keep the user name, but the user info is: %q", withBoth.GetUserInfo())
	}

	if withBoth.GetPassword() != "secret" {
		t.Errorf("GetPassword must be the password, but is: %q", withBoth.GetPassword())
	}
}

func TestAUriWithNoUserNameCarriesNoPassword(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeHttp, exampleHost, 0, "")

	withPassword := built.WithUserInfo("", "secret")

	if withPassword.GetPassword() != "" {
		t.Errorf("a URI with no user name must carry no password, but carries: %q", withPassword.GetPassword())
	}

	if withPassword.GetUserInfo() != "" {
		t.Errorf("the user info must be empty, but is: %q", withPassword.GetUserInfo())
	}
}

func TestEachWithMethodKeepsTheReceiverValueWhereTheNewOneIsInvalid(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeHttp, exampleHost, 8080, basicPath)

	if built.WithPort(70000).GetPort() != 8080 {
		t.Error("WithPort must keep the port where the new one is invalid, but did not")
	}

	if built.WithPath("/other?a=b").GetPath() != basicPath {
		t.Error("WithPath must keep the path where the new one carries a query string, but did not")
	}

	if built.WithQuery("a=b#part").GetQuery() != "" {
		t.Error("WithQuery must keep the query string where the new one carries a fragment, but did not")
	}
}

func TestGetSchemeAndGetUsernameReadWhatTheUriHolds(t *testing.T) {
	t.Parallel()

	built, err := uri.NewUri(constant.SchemeHttps, "user", "", exampleHost, 0, "", "", "")
	if err != nil {
		t.Fatalf("NewUri must build the URI, but reported: %v", err)
	}

	if built.GetScheme() != constant.SchemeHttps {
		t.Errorf("GetScheme must be the scheme, but is: %q", built.GetScheme())
	}

	if built.GetUsername() != "user" {
		t.Errorf("GetUsername must be the user name, but is: %q", built.GetUsername())
	}
}

func TestAnUnknownSchemeHasNoStandardPort(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.Scheme("ftp"), exampleHost, 0, "")

	if built.HasPort() {
		t.Error("an unknown scheme must name no port, but named one")
	}
}

func TestTheEmptySchemeReportsAStandardPortByItsHost(t *testing.T) {
	t.Parallel()

	withHost := newUri(t, constant.SchemeEmpty, exampleHost, 0, "")
	withNoHost := newUri(t, constant.SchemeEmpty, "", 0, "")

	if withHost.GetAuthority() != exampleHost {
		t.Errorf("the authority must be the host alone, but is: %q", withHost.GetAuthority())
	}

	if withNoHost.GetAuthority() != "" {
		t.Errorf("the authority must be empty with no host, but is: %q", withNoHost.GetAuthority())
	}
}

func TestTheEmptySchemeStatesAPortThatItNames(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeEmpty, exampleHost, 8080, "")

	if built.GetAuthority() != hostPort8080 {
		t.Errorf("the authority must state the port, but is: %q", built.GetAuthority())
	}
}

func TestTheUriSatisfiesItsContract(t *testing.T) {
	t.Parallel()

	var built contract.UriContract = newUri(t, constant.SchemeHttp, exampleHost, 0, basicPath)

	if built.GetPath() != basicPath {
		t.Errorf("the contract must read the path, but read: %q", built.GetPath())
	}
}

func TestAnUnknownSchemeWithAPortIsNeverStandard(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.Scheme("ftp"), exampleHost, 8080, "")

	if built.GetAuthority() != hostPort8080 {
		t.Errorf("an unknown scheme must state its port, but the authority is: %q", built.GetAuthority())
	}
}

func TestASchemeWithNoHostReportsAStandardPort(t *testing.T) {
	t.Parallel()

	built := newUri(t, constant.SchemeHttp, "", 8080, "")

	if built.GetPort() != 0 {
		t.Errorf("a URI with no host must report no port, but reports: %d", built.GetPort())
	}
}
