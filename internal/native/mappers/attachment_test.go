// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers_test

import (
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/misc"
)

// TestNullAttachment checks that the mappers of Message and Configuration turn a NULL attachment into a nil one: since
// SDK 3.18.0 the SDK passes NULL for an attachment of a class it does not deserialize.
func TestNullAttachment(t *testing.T) {
	available := mappers.AvailableMappers()
	for _, event := range []events.EventType{misc.NewMessage("MESSAGE"), misc.NewConfiguration("CONFIGURATION")} {
		mapper := available[event.Type().NativeCode()]
		// An event without an attachment gives a native event whose attachment is NULL.
		native := mapper.CEvent(event)
		back := mapper.GoEvent(native)
		mapper.FreeCEvent(native)
		var attachment *string
		switch e := back.(type) {
		case *misc.Message:
			attachment = e.Attachment()
		case *misc.Configuration:
			attachment = e.Attachment()
		default:
			t.Fatalf("GoEvent of %T = %T", event, back)
		}
		if attachment != nil {
			t.Errorf("%T: the attachment of a NULL attachment = %q; want nil", event, *attachment)
		}
	}
}
