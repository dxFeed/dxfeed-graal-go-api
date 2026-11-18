package events

import "strconv"

// EventFlag is a bit of the transactional event flags that indexed events carry; see IndexedEvent.
//
// The flags mark transactions (TxPending), removals (RemoveEvent) and snapshots (SnapshotBegin, SnapshotEnd,
// SnapshotSnip) in the stream of indexed events of one symbol and source. SnapshotMode and RemoveSymbol are used only
// inside the native layer and are never set on events that a subscription delivers.
type EventFlag int32

const (
	// TxPending marks a pending transaction: more events of the same transaction follow, and the state of the data
	// is consistent only after an event without this flag.
	TxPending EventFlag = 0x01
	// RemoveEvent marks the removal of the event with this index; the other fields of such an event carry no data.
	RemoveEvent EventFlag = 0x02
	// SnapshotBegin marks the first event of a snapshot.
	SnapshotBegin EventFlag = 0x04
	// SnapshotEnd marks the last event of a snapshot.
	SnapshotEnd EventFlag = 0x08
	// SnapshotSnip marks the last event of a snapshot that was cut at the time bound of the subscription.
	SnapshotSnip EventFlag = 0x10
	// SnapshotMode is a native-layer flag: the source sends a snapshot.
	SnapshotMode EventFlag = 0x40
	// RemoveSymbol is a native-layer flag: the symbol is removed from the subscription.
	RemoveSymbol EventFlag = 0x80
)

// Flag returns the flag as an int32 bit mask, for combining flags: TxPending.Flag() | SnapshotBegin.Flag().
func (f EventFlag) Flag() int32 {
	return int32(f)
}

// In reports whether the flag is set in the given event flags.
func (f EventFlag) In(eventFlags int32) bool {
	return eventFlags&int32(f) != 0
}

// Set returns the event flags with this flag set.
func (f EventFlag) Set(eventFlags int32) int32 {
	return eventFlags | int32(f)
}

// Clear returns the event flags with this flag cleared.
func (f EventFlag) Clear(eventFlags int32) int32 {
	return eventFlags &^ int32(f)
}

// String returns the name of the flag, such as "SNAPSHOT_BEGIN", or the hexadecimal value of an unknown bit.
func (f EventFlag) String() string {
	switch f {
	case TxPending:
		return "TX_PENDING"
	case RemoveEvent:
		return "REMOVE_EVENT"
	case SnapshotBegin:
		return "SNAPSHOT_BEGIN"
	case SnapshotEnd:
		return "SNAPSHOT_END"
	case SnapshotSnip:
		return "SNAPSHOT_SNIP"
	case SnapshotMode:
		return "SNAPSHOT_MODE"
	case RemoveSymbol:
		return "REMOVE_SYMBOL"
	}
	return "0x" + strconv.FormatUint(uint64(uint32(f)), 16)
}

// IndexedEvent is implemented by events that carry transactional event flags and an index that is unique per symbol
// and source: Order, TimeAndSale, Candle, Series, Greeks and the other indexed and time series events.
type IndexedEvent interface {
	EventType
	// EventFlags returns the transactional event flags, a combination of EventFlag bits.
	EventFlags() int32
	// SetEventFlags sets the transactional event flags.
	SetEventFlags(eventFlags int32)
	// Index returns the index of the event, unique per symbol and source.
	Index() int64
}

// IsSnapshotBegin reports whether the event starts a snapshot.
func IsSnapshotBegin(event IndexedEvent) bool {
	return SnapshotBegin.In(event.EventFlags())
}

// IsSnapshotEnd reports whether the event ends a snapshot.
func IsSnapshotEnd(event IndexedEvent) bool {
	return SnapshotEnd.In(event.EventFlags())
}

// IsSnapshotSnip reports whether the event ends a snapshot cut at the time bound of the subscription.
func IsSnapshotSnip(event IndexedEvent) bool {
	return SnapshotSnip.In(event.EventFlags())
}

// IsSnapshotEndOrSnip reports whether the event ends a snapshot, cut or not.
func IsSnapshotEndOrSnip(event IndexedEvent) bool {
	return IsSnapshotEnd(event) || IsSnapshotSnip(event)
}

// IsPending reports whether the event is part of a pending transaction.
func IsPending(event IndexedEvent) bool {
	return TxPending.In(event.EventFlags())
}

// IsRemove reports whether the event removes the event with its index.
func IsRemove(event IndexedEvent) bool {
	return RemoveEvent.In(event.EventFlags())
}
