package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// SetSystemProperty sets a Java system property of the SDK, for example "dxfeed.address".
func SetSystemProperty(key string, value string) error {
	return native.SetSystemProperty(key, value)
}

// GetSystemProperty returns a Java system property of the SDK; an unset property is an empty string.
func GetSystemProperty(key string) (string, error) {
	return native.GetSystemProperty(key)
}
