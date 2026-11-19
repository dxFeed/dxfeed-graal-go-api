// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package auth creates the authorization tokens of the dxFeed services.
//
// A token authorizes a connection of an endpoint ("demo.dxfeed.com:7300[login=entitle:<token>]") or a request of an
// Instrument Profile reader. It has a scheme ("Basic", "Bearer" or a custom one) and a value; its string form, the
// value of the HTTP Authorization header, is "<scheme> <value>":
//
//	token, err := auth.Parse("Bearer eyJhbGciOi...")
//	address := "demo.dxfeed.com:7300[login=entitle:" + token.Value() + "]"
//
// A [Token] is an immutable Go value; the SDK checks the arguments the same way as the Java API. Printing a token
// (fmt, log/slog) hides its value; [Token.Value] and [Token.HTTPAuthorization] give the secret.
package auth

import (
	"log/slog"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// The schemes of the tokens that the SDK creates.
const (
	// BasicScheme is the scheme of the user name and password tokens (RFC 2617).
	BasicScheme = "Basic"
	// BearerScheme is the scheme of the access tokens (RFC 6750).
	BearerScheme = "Bearer"
)

// Token is an authorization token. The zero value is not a token; use the functions of the package to create one.
type Token struct {
	scheme, value, httpAuthorization string
	user, password                   string
	basic                            bool
}

func newToken(kind int, args ...string) (*Token, error) {
	fields, err := native.NewAuthToken(kind, args...)
	if err != nil || fields == nil {
		return nil, err
	}
	return &Token{
		scheme:            fields.Scheme,
		value:             fields.Value,
		httpAuthorization: fields.HTTPAuthorization,
		user:              fields.User,
		password:          fields.Password,
		basic:             fields.HasUser,
	}, nil
}

// Parse parses the value of the HTTP Authorization header: "<scheme> <value>", for example "Bearer <token>" or
// "Basic <Base64 of user:password>". The scheme and the value are separated by a space.
func Parse(s string) (*Token, error) {
	return newToken(native.AuthTokenValueOf, s)
}

// NewBasic creates a Basic token from the string "user:password"; a string without a colon is an error.
func NewBasic(userPassword string) (*Token, error) {
	return newToken(native.AuthTokenBasic, userPassword)
}

// NewBasicUser creates a Basic token of the user name and the password.
func NewBasicUser(user, password string) (*Token, error) {
	return newToken(native.AuthTokenBasicUser, user, password)
}

// NewBasicOrNil is [NewBasicUser] that returns nil (and no error) when both the user name and the password are empty;
// check the result for nil before calling its methods.
func NewBasicOrNil(user, password string) (*Token, error) {
	return newToken(native.AuthTokenBasicOrNull, user, password)
}

// NewBearer creates a Bearer token; an empty token is an error.
func NewBearer(token string) (*Token, error) {
	return newToken(native.AuthTokenBearer, token)
}

// NewBearerOrNil is [NewBearer] that returns nil (and no error) for an empty token; check the result for nil before
// calling its methods.
func NewBearerOrNil(token string) (*Token, error) {
	return newToken(native.AuthTokenBearerOrNull, token)
}

// NewCustom creates a token of a custom scheme; an empty scheme or value is an error.
func NewCustom(scheme, value string) (*Token, error) {
	return newToken(native.AuthTokenCustom, scheme, value)
}

// Scheme returns the scheme of the token, for example [BasicScheme] or [BearerScheme].
func (t Token) Scheme() string { return t.scheme }

// Value returns the value of the token: the access token of a Bearer token, the Base64 of "user:password" of a Basic
// token.
func (t Token) Value() string { return t.value }

// HTTPAuthorization returns the value of the HTTP Authorization header, "<scheme> <value>". [Parse] of it gives an
// equal token.
func (t Token) HTTPAuthorization() string { return t.httpAuthorization }

// User returns the user name of a Basic token; false for the other schemes.
func (t Token) User() (string, bool) { return t.user, t.basic }

// Password returns the password of a Basic token; false for the other schemes.
func (t Token) Password() (string, bool) { return t.password, t.basic }

// String returns the scheme and a hidden value, "<scheme> ***", so that a logged token does not show the secret. Use
// [Token.HTTPAuthorization] for the whole value.
func (t Token) String() string {
	return t.scheme + " ***"
}

// GoString hides the value in the %#v format of fmt.
func (t Token) GoString() string {
	return "auth.Token{" + t.String() + "}"
}

// LogValue hides the value in log/slog.
func (t Token) LogValue() slog.Value {
	return slog.StringValue(t.String())
}
