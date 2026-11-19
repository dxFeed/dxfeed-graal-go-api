// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package auth_test

import (
	"fmt"
	"log"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/auth"
)

// A Basic token holds a user and a password; printing a token hides its value.
func ExampleNewBasicUser() {
	token, err := auth.NewBasicUser("demo", "secret")
	if err != nil {
		log.Fatal(err)
	}
	user, _ := token.User()
	fmt.Println(token.Scheme(), user)
	fmt.Println(token.HTTPAuthorization())
	fmt.Println(token)
	// Output:
	// Basic demo
	// Basic ZGVtbzpzZWNyZXQ=
	// Basic ***
}

// Parse reads a token from its HTTP Authorization form, "<scheme> <value>"; the value goes into the address of an
// endpoint.
func ExampleParse() {
	token, err := auth.Parse("Bearer example-token")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("demo.dxfeed.com:7300[login=entitle:" + token.Value() + "]")
	// Output: demo.dxfeed.com:7300[login=entitle:example-token]
}
