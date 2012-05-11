package xevent

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
)

// ClientMessageEvent embeds the struct by the same name from the xgb library.
type ClientMessageEvent struct {
	*xproto.ClientMessageEvent
}

// The unique code for a ClientMessage event.
const ClientMessage = xproto.ClientMessage

// NewClientMessage takes all arguments required to build a ClientMessageEvent 
// struct and hides the messy details.
// The varidic parameters coincide with the "data" part of a client message.
// Right now, this function only supports a list of up to 5 uint32s.
// XXX: Use type assertions to support bytes and uint16s.
func NewClientMessage(Format byte, Window xproto.Window, Type xproto.Atom,
	data ...interface{}) (*ClientMessageEvent, error) {

	// Create the client data list first
	var clientData xproto.ClientMessageDataUnion

	// Don't support formats 8 or 16 yet. They aren't used in EWMH anyway.
	switch Format {
	case 8:
		buf := make([]byte, 20)
		for i := 0; i < 20; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = data[i].(byte)
		}
		clientData = xproto.ClientMessageDataUnionData8New(buf)
	case 16:
		buf := make([]uint16, 10)
		for i := 0; i < 10; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = uint16(data[i].(int16))
		}
		clientData = xproto.ClientMessageDataUnionData16New(buf)
	case 32:
		buf := make([]uint32, 5)
		for i := 0; i < 5; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = uint32(data[i].(int))
		}
		clientData = xproto.ClientMessageDataUnionData32New(buf)
	default:
		return nil, fmt.Errorf("NewClientMessage: Unsupported format '%d'.",
			Format)
	}

	return &ClientMessageEvent{&xproto.ClientMessageEvent{
		Format: Format,
		Window: Window,
		Type:   Type,
		Data:   clientData,
	}}, nil
}

// ConfigureNotifyEvent embeds the struct by the same name in XGB.
type ConfigureNotifyEvent struct {
	*xproto.ConfigureNotifyEvent
}

// The unique code for a ConfigureNotify event.
const ConfigureNotify = xproto.ConfigureNotify

// NewConfigureNotify takes all arguments required to build a 
// ConfigureNotifyEvent struct and hides the messy details.
func NewConfigureNotify(Event, Window, AboveSibling xproto.Window,
	X, Y, Width, Height int, BorderWidth uint16,
	OverrideRedirect bool) *ConfigureNotifyEvent {

	return &ConfigureNotifyEvent{&xproto.ConfigureNotifyEvent{
		Event: Event, Window: Window, AboveSibling: AboveSibling,
		X: int16(X), Y: int16(Y), Width: uint16(Width), Height: uint16(Height),
		BorderWidth: BorderWidth, OverrideRedirect: OverrideRedirect,
	}}
}
