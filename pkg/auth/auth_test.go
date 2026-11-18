package auth_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/auth"
)

func TestBasicRoundTrip(t *testing.T) {
	token, err := auth.NewBasicUser("user", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if token.Scheme() != auth.BasicScheme || token.Value() != "dXNlcjpwYXNz" ||
		token.HTTPAuthorization() != "Basic dXNlcjpwYXNz" {
		t.Errorf("token = %q %q %q", token.Scheme(), token.Value(), token.HTTPAuthorization())
	}
	parsed, err := auth.Parse(token.HTTPAuthorization())
	if err != nil {
		t.Fatal(err)
	}
	if *parsed != *token {
		t.Errorf("parsed = %+v, want %+v", parsed, token)
	}
	if user, ok := parsed.User(); !ok || user != "user" {
		t.Errorf("User() = %q, %v", user, ok)
	}
	if password, ok := parsed.Password(); !ok || password != "pass" {
		t.Errorf("Password() = %q, %v", password, ok)
	}
	if fromString, err := auth.NewBasic("user:pass"); err != nil || *fromString != *token {
		t.Errorf("NewBasic = %+v, %v", fromString, err)
	}
}

func TestBearerAndCustom(t *testing.T) {
	bearer, err := auth.NewBearer("abc.def")
	if err != nil {
		t.Fatal(err)
	}
	if bearer.HTTPAuthorization() != "Bearer abc.def" || bearer.Value() != "abc.def" {
		t.Errorf("bearer = %q", bearer.HTTPAuthorization())
	}
	if _, ok := bearer.User(); ok {
		t.Error("a Bearer token has a user")
	}
	custom, err := auth.Parse("Custom x y")
	if err != nil {
		t.Fatal(err)
	}
	if custom.Scheme() != "Custom" || custom.Value() != "x y" {
		t.Errorf("custom = %q %q", custom.Scheme(), custom.Value())
	}
	if made, err := auth.NewCustom("Custom", "x y"); err != nil || *made != *custom {
		t.Errorf("NewCustom = %+v, %v", made, err)
	}
}

func TestOrNil(t *testing.T) {
	if token, err := auth.NewBasicOrNil("", ""); token != nil || err != nil {
		t.Errorf("NewBasicOrNil of empty = %v, %v", token, err)
	}
	if token, err := auth.NewBasicOrNil("u", ""); token == nil || err != nil {
		t.Errorf("NewBasicOrNil of a user = %v, %v", token, err)
	}
	if token, err := auth.NewBearerOrNil(""); token != nil || err != nil {
		t.Errorf("NewBearerOrNil of empty = %v, %v", token, err)
	}
}

func TestInvalid(t *testing.T) {
	for name, create := range map[string]func() (*auth.Token, error){
		"Parse without a space": func() (*auth.Token, error) { return auth.Parse("user:pass") },
		"Parse of empty":        func() (*auth.Token, error) { return auth.Parse("") },
		"NewBasic without ':'":  func() (*auth.Token, error) { return auth.NewBasic("userpass") },
		"NewBearer of empty":    func() (*auth.Token, error) { return auth.NewBearer("") },
		"NewCustom of empty":    func() (*auth.Token, error) { return auth.NewCustom("", "v") },
	} {
		if token, err := create(); token != nil || err == nil {
			t.Errorf("%s = %v, %v; want an error", name, token, err)
		}
	}
}

func TestStringHidesTheSecret(t *testing.T) {
	token, err := auth.NewBearer("secret")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{token.String(), fmt.Sprintf("%+v", token)} {
		if strings.Contains(s, "secret") || !strings.HasPrefix(s, "Bearer") {
			t.Errorf("string = %q", s)
		}
	}
	var none *auth.Token
	if none.String() != "<nil>" {
		t.Errorf("nil string = %q", none.String())
	}
}
