// Package apibridge gives the other packages of the module the native handles behind the types of pkg/api without
// exporting them from pkg/api. pkg/api sets the functions in its init, so a package that imports pkg/api can call
// them.
package apibridge

import "github.com/dxfeed/dxfeed-graal-go-api/internal/native"

var (
	// FeedHandle returns the native feed of a *api.DXFeed, or nil for a nil feed.
	FeedHandle func(feed any) *native.DXFeedHandle
	// ExecutorHandle returns the native executor of an api.Executor, or nil for a nil or closed executor.
	ExecutorHandle func(executor any) *native.ExecutorHandle
)
