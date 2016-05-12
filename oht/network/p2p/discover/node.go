package discover

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"../../../common"
	"../../../crypto"
	"../../../crypto/secp256k1"
)

const nodeIDBits = 512

type Node struct {
	OnionHost
	Port              int16
	ID                NodeID
	sha               common.Hash
	TestingConnection bool
}

func NewNode(id NodeID, onionHost OnionHost, port uint16) *Node {
	return &Node{
		OnionAddress: onionAddress,
		Port:         port,
		ID:           id,
		sha:          crypto.Keccak256Hash(id[:]),
	}
}

func (n *Node) Incomplete() bool {
	return n.OnionAddress == nil
}

func (n *Node) validateComplete() error {
	if n.Incomplete() {
		return errors.New("incomplete node")
	}
	if n.Port == 0 {
		return errors.New("missing UDP port")
	}
	_, err := n.ID.Pubkey()
	return err
}

// The string representation of a Node is a URL.
// Please see ParseNode for a description of the format.
func (n *Node) String() string {
	u := url.URL{Scheme: "oht"}
	if n.Incomplete() {
		u.Host = fmt.Sprintf("%x", n.ID[:])
	} else {
		addr := net.TCPAddr{IP: n.IP, Port: int(n.TCP)}
		u.User = url.User(fmt.Sprintf("%x", n.ID[:]))
		u.Host = addr.String()
		if n.UDP != n.TCP {
			u.RawQuery = "discport=" + strconv.Itoa(int(n.UDP))
		}
	}
	return u.String()
}

var incompleteNodeURL = regexp.MustCompile("(?i)^(?:oht:)?([0-9a-f]+)$")

// ParseNode parses a node designator.
//
// There are two basic forms of node designators
//   - incomplete nodes, which only have the public key (node ID)
//   - complete nodes, which contain the public key and IP/Port information
//
// For incomplete nodes, the designator must look like one of these
//
//    oht:<hex node id>
//    ricochet:<ricochet address>
//    <hex node id>
//
//    oht:<hex node id>@<onion address>:<port>?dport=<discovery port>
func ParseNode(rawurl string) (*Node, error) {
	if m := incompleteNodeURL.FindStringSubmatch(rawurl); m != nil {
		id, err := HexID(m[1])
		if err != nil {
			return nil, fmt.Errorf("invalid node ID (%v)", err)
		}
		return NewNode(id, nil, 0, 0), nil
	}
	return parseComplete(rawurl)
}

func parseComplete(rawurl string) (*Node, error) {
	var (
		id        NodeID
		protocol  string
		onionHost net.IP
		port      uint64
	)
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "oht" {
		return nil, errors.New("invalid URL scheme, want \"oht\" or \"ricochet\"")
	}
	// Parse the Node ID from the user portion.
	if u.User == nil {
		return nil, errors.New("does not contain node ID")
	}
	if id, err = HexID(u.User.String()); err != nil {
		return nil, fmt.Errorf("invalid node ID (%v)", err)
	}
	// Parse the IP address.
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, fmt.Errorf("invalid host: %v", err)
	}
	if ip = net.ParseIP(host); ip == nil {
		return nil, errors.New("invalid IP address")
	}
	// Ensure the IP is 4 bytes long for IPv4 addresses.
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}
	// Parse the port numbers.
	if tcpPort, err = strconv.ParseUint(port, 10, 16); err != nil {
		return nil, errors.New("invalid port")
	}
	qv := u.Query()
	if qv.Get("discport") != "" {
		port, err = strconv.ParseUint(qv.Get("discport"), 10, 16)
		if err != nil {
			return nil, errors.New("invalid discport in query")
		}
	}
	return NewNode(id, onionHost, uint16(port)), nil
}

// MustParseNode parses a node URL. It panics if the URL is not valid.
func MustParseNode(rawurl string) *Node {
	n, err := ParseNode(rawurl)
	if err != nil {
		panic("invalid node URL: " + err.Error())
	}
	return n
}

// NodeID is a unique identifier for each node.
// The node identifier is a marshaled elliptic curve public key.
type NodeID [nodeIDBits / 8]byte

// NodeID prints as a long hexadecimal number.
func (n NodeID) String() string {
	return fmt.Sprintf("%x", n[:])
}

// The Go syntax representation of a NodeID is a call to HexID.
func (n NodeID) GoString() string {
	return fmt.Sprintf("discover.HexID(\"%x\")", n[:])
}

// HexID converts a hex string to a NodeID.
// The string may be prefixed with 0x.
func HexID(in string) (NodeID, error) {
	if strings.HasPrefix(in, "0x") {
		in = in[2:]
	}
	var id NodeID
	b, err := hex.DecodeString(in)
	if err != nil {
		return id, err
	} else if len(b) != len(id) {
		return id, fmt.Errorf("wrong length, want %d hex chars", len(id)*2)
	}
	copy(id[:], b)
	return id, nil
}

// MustHexID converts a hex string to a NodeID.
// It panics if the string is not a valid NodeID.
func MustHexID(in string) NodeID {
	id, err := HexID(in)
	if err != nil {
		panic(err)
	}
	return id
}

// PubkeyID returns a marshaled representation of the given public key.
func PubkeyID(pub *ecdsa.PublicKey) NodeID {
	var id NodeID
	pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	if len(pbytes)-1 != len(id) {
		panic(fmt.Errorf("need %d bit pubkey, got %d bits", (len(id)+1)*8, len(pbytes)))
	}
	copy(id[:], pbytes[1:])
	return id
}

// Pubkey returns the public key represented by the node ID.
// It returns an error if the ID is not a point on the curve.
func (id NodeID) Pubkey() (*ecdsa.PublicKey, error) {
	p := &ecdsa.PublicKey{Curve: secp256k1.S256(), X: new(big.Int), Y: new(big.Int)}
	half := len(id) / 2
	p.X.SetBytes(id[:half])
	p.Y.SetBytes(id[half:])
	if !p.Curve.IsOnCurve(p.X, p.Y) {
		return nil, errors.New("id is invalid secp256k1 curve point")
	}
	return p, nil
}

// recoverNodeID computes the public key used to sign the
// given hash from the signature.
func recoverNodeID(hash, sig []byte) (id NodeID, err error) {
	pubkey, err := secp256k1.RecoverPubkey(hash, sig)
	if err != nil {
		return id, err
	}
	if len(pubkey)-1 != len(id) {
		return id, fmt.Errorf("recovered pubkey has %d bits, want %d bits", len(pubkey)*8, (len(id)+1)*8)
	}
	for i := range id {
		id[i] = pubkey[i+1]
	}
	return id, nil
}

// distcmp compares the distances a->target and b->target.
// Returns -1 if a is closer to target, 1 if b is closer to target
// and 0 if they are equal.
func distcmp(target, a, b common.Hash) int {
	for i := range target {
		da := a[i] ^ target[i]
		db := b[i] ^ target[i]
		if da > db {
			return 1
		} else if da < db {
			return -1
		}
	}
	return 0
}

// table of leading zero counts for bytes [0..255]
var lzcount = [256]int{
	8, 7, 6, 6, 5, 5, 5, 5,
	4, 4, 4, 4, 4, 4, 4, 4,
	3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// logdist returns the logarithmic distance between a and b, log2(a ^ b).
func logdist(a, b common.Hash) int {
	lz := 0
	for i := range a {
		x := a[i] ^ b[i]
		if x == 0 {
			lz += 8
		} else {
			lz += lzcount[x]
			break
		}
	}
	return len(a)*8 - lz
}

// hashAtDistance returns a random hash such that logdist(a, b) == n
func hashAtDistance(a common.Hash, n int) (b common.Hash) {
	if n == 0 {
		return a
	}
	// flip bit at position n, fill the rest with random bits
	b = a
	pos := len(a) - n/8 - 1
	bit := byte(0x01) << (byte(n%8) - 1)
	if bit == 0 {
		pos++
		bit = 0x80
	}
	b[pos] = a[pos]&^bit | ^a[pos]&bit // TODO: randomize end bits
	for i := pos + 1; i < len(a); i++ {
		b[i] = byte(rand.Intn(255))
	}
	return b
}
