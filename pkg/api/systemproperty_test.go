package api

import (
	"fmt"
	"sync"
	"testing"
)

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

func TestGetSystemPropertyWithUnicode(t *testing.T) {
	key := "∑🍕🍔🍗"
	value := "∏🍻🍺"
	SetSystemProperty(key, value)
	returnValue := GetSystemProperty(key)
	if returnValue != value {
		t.Fatalf(`GetSystemProperty("%s") should be "%s". But it equals "%s"`, key, value, returnValue)
	}
}

func TestConcurrentGetSetSystemProperty(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		num := i
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("test_key_%d", num)
			value := fmt.Sprintf("test_value_%d", num)
			SetSystemProperty(key, value)
			returnValue := GetSystemProperty(key)
			if returnValue != value {
				t.Errorf(`GetSystemProperty("%s") should be "%s". But it equals "%s"`, key, value, returnValue)
				return
			}
		}()
	}
	wg.Wait()
}
