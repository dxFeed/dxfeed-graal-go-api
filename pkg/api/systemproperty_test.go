package api

import "testing"

func TestGetSystemProperty(t *testing.T) {
	key := "test_key"
	value := "test_value"
	SetSystemProperty(key, value)
	returnValue := GetSystemProperty(key)
	if returnValue != value {
		t.Fatalf(`GetSystemProperty("%s") should be "%s". But it equals "%s"`, key, value, returnValue)
	}
}

func TestGetEmptySystemProperty(t *testing.T) {
	key := "test_key"
	value := ""
	SetSystemProperty(key, value)
	returnValue := GetSystemProperty(key)
	if returnValue != value {
		t.Fatalf(`GetSystemProperty("%s") should be "%s". But it equals "%s"`, key, value, returnValue)
	}
}
