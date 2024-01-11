package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

func SetSystemProperty(key string, value string) {
	native.SetSystemProperty(key, value)
}

func GetSystemProperty(key string) string {
	return native.GetSystemProperty(key)
}
