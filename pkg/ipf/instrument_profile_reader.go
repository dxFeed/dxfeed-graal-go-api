package ipf

import "C"
import (
	"io"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/auth"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// The reader owns a native object: close it.
var _ io.Closer = (*InstrumentProfileReader)(nil)

// InstrumentProfileReader reads instrument profiles; close it when it is no longer needed.
type InstrumentProfileReader struct {
	reader *native.InstrumentProfileReader
}

// NewInstrumentProfileReader creates a reader.
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

// ResolveSourceURL turns an address into the URL to read profiles from; a plain host:port becomes a URL of the
// instrument profile service.
func ResolveSourceURL(address string) (*string, error) {
	return native.ResolveSourceURL(address)
}

// Close releases the native resources of the reader.
func (r *InstrumentProfileReader) Close() error {
	return r.reader.Close()
}

// LastModified returns the modification time of the last read file, in milliseconds since the Unix epoch, or 0 if it
// is unknown.
func (r *InstrumentProfileReader) LastModified() (int64, error) {
	return r.reader.GetLastModified()
}

// WasComplete reports whether the last read file was read completely.
func (r *InstrumentProfileReader) WasComplete() (bool, error) {
	return r.reader.WasComplete()
}

// ReadFromFile reads the profiles from a file path or URL; a .zip or .gz file is unpacked.
func (r *InstrumentProfileReader) ReadFromFile(address string) ([]*events.InstrumentProfile, error) {
	return r.reader.ReadFromFile(address)
}

// ReadFromFileWithPassword reads the profiles from a URL that needs a user and a password.
func (r *InstrumentProfileReader) ReadFromFileWithPassword(address string, user string, password string) ([]*events.InstrumentProfile, error) {
	return r.reader.ReadFromFileWithPassword(address, user, password)
}

// ReadFromFileWithToken reads the profiles from a URL that needs an authorization token; a nil token reads without
// authorization, like ReadFromFile.
func (r *InstrumentProfileReader) ReadFromFileWithToken(address string, token *auth.Token) ([]*events.InstrumentProfile, error) {
	if token == nil {
		return r.reader.ReadFromFile(address)
	}
	return r.reader.ReadFromFileWithToken(address, token.Scheme(), token.Value())
}

// WriteToFile writes the profiles to a file in the IPF format; a name ending with .gz or .zip compresses the file.
// An existing file is replaced, and an I/O error may leave an incomplete file. The profiles are not validated beyond
// a nil profile, a NUL character and the int32 range of ICB, SIC and the dates. The file can be read back with
// InstrumentProfileReader.ReadFromFile.
func WriteToFile(path string, profiles []*events.InstrumentProfile) error {
	return native.WriteProfilesToFile(path, profiles)
}
