package dxendpoint

type DXEndpoint struct {
}

func Create() DXEndpoint {
	return DXEndpoint{}
}

func (e DXEndpoint) Connect(address string) {
}

func (e DXEndpoint) Close() {

}
