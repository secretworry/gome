package net

type (
	FrameType uint8
	Flags uint8
)


// A FrameHeader is a 9 byte header of our protocol
type FrameHeader struct {
	// Type is 1 byte frame type
	Type FrameType
	// Flags are 1 byte of 8 potential bit flags per frame.
	// They are specified to the frame type
	Flags Flags
	// Length is the length of the frame, not including the 9 byte header.
	// The maximum size is 2^24 - 1( 24KB), but due to the MTU, no more than 1460b is suggested
	Length uint32
	// SrcChannelID is which channel in the sender the frame is for. For non-channel-specified the field is 0,
	// otherwise the sender should fill in the channel it sources from.
	SrcChannelID uint32
	// DestChannelID is which channel in the receiver the frame is for. When sending the first piece of data, or a
	// non-channel-specified frame, the field is 0
	DestChannelID uint32
}

const (
	FrameData	FrameType = 0x0
	FramePing	FrameType = 0x1
	FrameGoAway	FrameType = 0x2
)

const (
	// Data Frame
	FlagDataEndStream	Flags = 0x1

	// Ping Frame
	FlagPingAck		Flags = 0x1
)

var flagName = map[FrameType]map[Flags]string{
	FrameData: {
		FlagDataEndStream: "END_STREAM",
	},
	FramePing: {
		FlagPingAck: "ACK",
	},
}

const (
	frameHeaderLength = 13
)

