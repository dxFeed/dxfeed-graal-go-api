package api

type DXEndpoint struct {
	role Role
}

type Role int32

const (
	Feed       Role = iota
	StreamFeed Role = iota
)

func CreateEndpoint(role Role) DXEndpoint {
	endpoint := DXEndpoint{role: role}
	return endpoint
}

func (e *DXEndpoint) Connect(address string) {
}

func (e *DXEndpoint) Close() {
}

func (e *DXEndpoint) GetFeed() DXFeed {
	return DXFeed{}
}
