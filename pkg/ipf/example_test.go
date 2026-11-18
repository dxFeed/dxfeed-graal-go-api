package ipf_test

import (
	"errors"
	"fmt"
	"log"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
)

// A Java exception inside the SDK, here for a missing file, is returned as common.JavaError.
func ExampleInstrumentProfileReader_ReadFromFile() {
	reader, err := ipf.NewInstrumentProfileReader()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = reader.Close() }()

	_, err = reader.ReadFromFile("no-such-file.ipf")
	var javaErr common.JavaError
	fmt.Println(errors.As(err, &javaErr), javaErr.ClassName != "")
	// Output: true true
}
