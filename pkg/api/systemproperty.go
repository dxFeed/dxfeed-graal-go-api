package api

import "dxfeed-graal-go-api/pkg/native"

func SetSystemProperty(key string, value string) {
	native.SetSystemProperty(key, value)
}

func GetSystemProperty(key string) string {
	return native.GetSystemProperty(key)
}
