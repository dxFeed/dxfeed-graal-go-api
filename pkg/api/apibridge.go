// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/apibridge"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// init gives the other packages of the module the native handles behind DXFeed and Executor.
func init() {
	apibridge.FeedHandle = func(feed any) *native.DXFeedHandle {
		f, ok := feed.(*DXFeed)
		if !ok || f == nil {
			return nil
		}
		return f.feed
	}
	apibridge.ExecutorHandle = func(executor any) *native.ExecutorHandle {
		e, ok := executor.(Executor)
		if !ok || e == nil {
			return nil
		}
		h := e.handle()
		if h == nil || h.IsFreed() {
			return nil
		}
		return h
	}
}
