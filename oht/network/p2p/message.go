package network

import (
	"crypto/sha512"
	"encoding/binary"
	"github.com/egordon/gobitmsg/types"
	"strconv"
)

var ProtocolName = "oht"

type errCode int

const (
	NetworkId          = 1
	ProtocolMaxMsgSize = 10 * 1024 * 1024 // Maximum cap on the size of a protocol message
	ErrMsgTooLarge     = iota
	ErrDecode
	ErrInvalidMsgCode
	ErrProtocolVersionMismatch
	ErrNetworkIdMismatch
	ErrGenesisBlockMismatch
	ErrNoStatusMsg
	ErrExtraStatusMsg
	ErrSuspendedPeer
	headerLen  = 24
	knownMagic = 0xe9beb4d9
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

var errorToString = map[int]string{
	ErrMsgTooLarge:             "Message too long",
	ErrDecode:                  "Invalid message",
	ErrInvalidMsgCode:          "Invalid message code",
	ErrProtocolVersionMismatch: "Protocol version mismatch",
	ErrNetworkIdMismatch:       "NetworkId mismatch",
	ErrGenesisBlockMismatch:    "Genesis block mismatch",
	ErrNoStatusMsg:             "No status message",
	ErrExtraStatusMsg:          "Extra status message",
	ErrSuspendedPeer:           "Suspended peer",
}

// Message holds a standard, serializeable BitMsg message header and generic payload.
// See more info at https://bitmessage.org/wiki/Protocol_specification#Message_structure
type Message struct {
	peer     *Peer  // Peer who sent the message (nil if local is sending)
	magic    uint32 // Magic number associated with network
	command  string // Action this message wants to take
	length   uint32 // Length of the payload
	checksum uint32 // First 4 bytes of sha512(payload)
	payload  []byte
}

// statusData is the network packet for the status message.
type statusData struct {
	ProtocolVersion uint32
	NetworkId       uint32
	TD              *big.Int
	CurrentBlock    common.Hash
	GenesisBlock    common.Hash
}

// MakeMessage generates a new message given a command, payload, and recipient.
// No defensive copy is made of the byte slice
func MakeMessage(cmd string, pload types.Serializer, recipient *Peer) *Message {
	msg := new(Message)
	msg.peer = recipient
	msg.magic = knownMagic
	msg.payload = pload.Serialize()
	msg.length = uint32(len(msg.payload))
	msg.command = strconv.Quote(cmd)

	digest := sha512.Sum512(msg.payload)
	checksum := make([]byte, 4, 4)
	for i := 0; i < 4; i++ {
		checksum[i] = digest[i]
	}
	msg.checksum = binary.BigEndian.Uint32(checksum)
	return msg
}

// Serialize converts a message to a byte stream that can be sent over the network.
func (m *Message) Serialize() []byte {
	ret := make([]byte, headerLen, headerLen)
	binary.BigEndian.PutUint32(ret[0:4], m.magic)
	copy(ret[4:16], m.command)

	// Ensure string had a max size of 12
	binary.BigEndian.PutUint32(ret[16:20], m.length)
	binary.BigEndian.PutUint32(ret[20:24], m.checksum)
	ret = append(ret, m.payload...)

	return ret
}

func (m *Message) validate() error {

	if m.magic != knownMagic {
		return MessageError(EMAGIC)
	}

	if m.length != uint32(len(m.payload)) {
		return MessageError(EPALEN)
	}

	checksum := make([]byte, 4, 4)
	digest := sha512.Sum512(m.payload)

	checksum[0] = digest[0]
	checksum[1] = digest[1]
	checksum[2] = digest[2]
	checksum[3] = digest[3]

	if m.checksum != binary.BigEndian.Uint32(checksum) {
		return MessageError(ECHECK)
	}

	return nil
}

func (m *Message) makeHeader(rawBytes []byte) error {
	if len(rawBytes) < headerLen {
		return MessageError(ESMALL)
	}
	m.magic = binary.BigEndian.Uint32(rawBytes[:4])
	m.command = string(rawBytes[4:16])
	m.length = binary.BigEndian.Uint32(rawBytes[16:20])
	m.checksum = binary.BigEndian.Uint32(rawBytes[20:24])
	m.payload = nil
	return nil
}

func (m *Message) setPayload(rawBytes []byte) {
	m.payload = make([]byte, len(rawBytes), len(rawBytes))
	copy(m.payload, rawBytes)
}

// Payload returns the generic payload byte slice of the Message
func (m *Message) Payload() []byte {
	if m == nil {
		return nil
	} else {
		return m.payload
	}
}

// Command returns the command associated with the message
func (m *Message) Command() string {
	if m == nil {
		return ""
	} else {
		return m.command
	}
}

// Sender returns the recipient and/or sender of the message
func (m *Message) Sender() *Peer {
	if m == nil {
		return nil
	} else {
		return m.peer
	}
}
