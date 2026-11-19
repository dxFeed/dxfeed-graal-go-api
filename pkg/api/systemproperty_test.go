// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"sync"
	"testing"
)

func TestGetSystemProperty(t *testing.T) {
	key := "test_key"
	value := "test_value"
	if err := SetSystemProperty(key, value); err != nil {
		t.Fatal(err)
	}
	returnValue, err := GetSystemProperty(key)
	if err != nil {
		t.Fatal(err)
	}
	if returnValue != value {
		t.Fatalf(`GetSystemProperty("%s") should be "%s". But it equals "%s"`, key, value, returnValue)
	}
}

func TestGetEmptySystemProperty(t *testing.T) {
	key := "test_key"
	value := ""
	if err := SetSystemProperty(key, value); err != nil {
		t.Fatal(err)
	}
	returnValue, err := GetSystemProperty(key)
	if err != nil {
		t.Fatal(err)
	}
	if returnValue != value {
		t.Fatalf(`GetSystemProperty("%s") should be "%s". But it equals "%s"`, key, value, returnValue)
	}
}

func TestGetSystemPropertyWithUnicode(t *testing.T) {
	key := "∑🍕🍔🍗"
	value := "∏🍻🍺"
	if err := SetSystemProperty(key, value); err != nil {
		t.Fatal(err)
	}
	returnValue, err := GetSystemProperty(key)
	if err != nil {
		t.Fatal(err)
	}
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
			if err := SetSystemProperty(key, value); err != nil {
				t.Error(err)
				return
			}
			returnValue, err := GetSystemProperty(key)
			if err != nil {
				t.Error(err)
				return
			}
			if returnValue != value {
				t.Errorf(`GetSystemProperty("%s") should be "%s". But it equals "%s"`, key, value, returnValue)
				return
			}
		}()
	}
	wg.Wait()
}
