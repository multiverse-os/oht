package network

// PeerError handles all errors associated with the Peer struct.
type PeerError int

// Error constants for PeerError
const (
	ENTFND = iota
	EPRDUP = iota
)

// Error converts the PeerError constant to a human-readable string
func (errno PeerError) Error() string {
	switch errno {
	case ENTFND:
		return "Peer not in Peer List"
	case EPRDUP:
		return "Peer already in list and connected"
	default:
		return "Unkown Peer Error"
	}
}

// MessageError handles all errors associated with the Message struct.
type MessageError int

// Error constants for MessageError
const (
	ESMALL = iota
	EMAGIC = iota
	EPALEN = iota
	ECHECK = iota
	ETIMEO = iota
)

// Error converts the MessageError constant to a human-readable string
func (errno MessageError) Error() string {
	switch errno {
	case ESMALL:
		return "Header size too small"
	case EMAGIC:
		return "Magic number in header is unknown"
	case EPALEN:
		return "Payload length does not match header"
	case ECHECK:
		return "Invalid SHA512 checksum for payload"
	case ETIMEO:
		return "Message timeout, message queue empty"
	default:
		return "Unknown Message Error"
	}
}
