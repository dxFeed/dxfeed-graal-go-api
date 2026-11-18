package api

import (
	"errors"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

// EndpointBuilder configures and builds endpoints, as DXEndpoint.Builder of the dxFeed Java API. Create it with
// NewEndpointBuilder; the With methods return the builder, so the calls chain:
//
//	endpoint, err := api.NewEndpointBuilder().
//		WithRole(api.Feed).
//		WithName("quotes").
//		WithProperty("dxfeed.aggregationPeriod", "1s").
//		Build()
//
// The settings are applied in the order of the calls, so a later one wins over an earlier one with the same key.
// A builder holds no native resources and may be reused: every Build creates a new endpoint. It is not safe for
// concurrent use.
type EndpointBuilder struct {
	role     common.Role
	steps    []native.BuilderStep
	executor Executor
}

// NewEndpointBuilder creates a builder of endpoints with the Feed role.
func NewEndpointBuilder() *EndpointBuilder {
	return &EndpointBuilder{role: Feed}
}

// WithRole sets the role of the endpoints; the default is Feed.
func (b *EndpointBuilder) WithRole(role common.Role) *EndpointBuilder {
	b.role = role
	return b
}

// WithName sets the name of the endpoints, which tells several endpoints of one process apart in the logs. It is the
// same as the property "name".
func (b *EndpointBuilder) WithName(name string) *EndpointBuilder {
	b.steps = append(b.steps, native.WithName(name))
	return b
}

// WithProperty sets a property of the endpoints, such as "dxfeed.address"; it wins over the same property of a
// properties file and of the system properties. A property that the builder does not support is ignored: see
// SupportsProperty.
func (b *EndpointBuilder) WithProperty(key, value string) *EndpointBuilder {
	b.steps = append(b.steps, native.WithProperty(key, value))
	return b
}

// WithProperties sets the properties of the endpoints, as WithProperty for each of them.
func (b *EndpointBuilder) WithProperties(properties map[string]string) *EndpointBuilder {
	for key, value := range properties {
		b.WithProperty(key, value)
	}
	return b
}

// WithPropertiesFile sets the properties of the endpoints from a file in the format of Java properties files; the
// file is read by Build, which fails if it cannot be read.
func (b *EndpointBuilder) WithPropertiesFile(path string) *EndpointBuilder {
	b.steps = append(b.steps, native.WithPropertiesFile(path))
	return b
}

// WithExecutor sets the executor of the notifications of the endpoints: event listeners, state changes and the like.
// nil restores the SDK's default executor. Build fails if the executor is closed.
func (b *EndpointBuilder) WithExecutor(executor Executor) *EndpointBuilder {
	b.executor = executor
	return b
}

// SupportsProperty reports whether the builder supports a property; the unsupported ones are ignored.
func (b *EndpointBuilder) SupportsProperty(key string) (bool, error) {
	return native.BuilderSupportsProperty(key, native.WithRole(b.role))
}

// Build creates a new endpoint with the builder's settings. An endpoint with the role Feed or Publisher connects at
// once if the property "dxfeed.address" (Feed) or "dxpublisher.address" (Publisher) is set.
func (b *EndpointBuilder) Build() (*DXEndpoint, error) {
	steps := make([]native.BuilderStep, 0, len(b.steps)+2)
	steps = append(steps, native.WithRole(b.role))
	if b.executor != nil {
		h := b.executor.handle()
		if h.IsFreed() {
			return nil, errors.New("dxfeed: the executor of the endpoint builder is closed")
		}
		steps = append(steps, native.WithExecutor(h))
	}
	steps = append(steps, b.steps...)
	handle, err := native.BuildDXEndpointHandle(steps...)
	return newEndpoint(b.role, handle, err)
}
