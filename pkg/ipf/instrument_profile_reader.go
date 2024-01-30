package ipf

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

type InstrumentProfileReader struct {
	reader *native.InstrumentProfileReader
}

func NewInstrumentProfileReader() (*InstrumentProfileReader, error) {
	handle, err := native.NewInstrumentProfileReader()
	if err != nil {
		return nil, err
	}

	e := &InstrumentProfileReader{
		reader: handle,
	}
	return e, nil
}

func (r *InstrumentProfileReader) GetLastModified() (int64, error) {
	return r.reader.GetLastModified()
}

func (r *InstrumentProfileReader) WasComplete() (bool, error) {
	return r.reader.WasComplete()
}

func (r *InstrumentProfileReader) ReadFromFile(address string) ([]*events.InstrumentProfile, error) {
	return r.reader.ReadFromFile(address)
}

func (r *InstrumentProfileReader) ReadFromFileWithPassword(address string, user string, password string) ([]*events.InstrumentProfile, error) {
	return r.reader.ReadFromFileWithPassword(address, user, password)
}

func (r *InstrumentProfileReader) ReadDataWithAddress(data []byte, address string) ([]*events.InstrumentProfile, error) {
	return r.reader.ReadDataWithAddress(data, address)
}

func (r *InstrumentProfileReader) ReadData(data []byte) ([]*events.InstrumentProfile, error) {
	return r.reader.ReadData(data)
}

func ResolveSourceURL(address string) (*string, error) {
	return native.ResolveSourceURL(address)
}
