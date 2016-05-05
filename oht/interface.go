package oht

import (
	"log"
	"os"

	"./common"
	"./crypto"
	"./network"
	"./network/p2p"
	"./network/webui"
	"./types"
)

type Interface struct {
	config      *Config
	tor         *network.TorProcess
	webUIServer *webui.Server
	p2p         *p2p.Manager
}

func NewInterface(c *Config, t *network.TorProcess, w *webui.Server, p *p2p.Manager) (i *Interface) {
	return &Interface{
		config:      c,
		tor:         t,
		webUIServer: w,
		p2p:         p,
	}
}

// msg validator
type Validator interface{}

// GENERAL INFORMATION
func (i *Interface) ClientName() string    { return i.config.ClientName }
func (i *Interface) ClientVersion() string { return i.config.clientVersion() }
func (i *Interface) ClientInfo() string    { return i.config.clientInfo() }

// TOR INFORMATION
func (i *Interface) TorListenPort() string     { return i.tor.ListenPort }
func (i *Interface) TorSocksPort() string      { return i.tor.SocksPort }
func (i *Interface) TorControlPort() string    { return i.tor.ControlPort }
func (i *Interface) TorWebUIPort() string      { return i.tor.WebUIPort }
func (i *Interface) TorOnionHost() string      { return i.tor.OnionHost }
func (i *Interface) TorWebUIOnionHost() string { return i.tor.WebUIOnionHost }

func (i *Interface) ProtocolVersion() {}
func (i *Interface) Locale() string   { return "en" }
func (i *Interface) PeerCount() int   { return i.PeerCount() }
func (i *Interface) MaxPeers() int    { return i.MaxPeers() }

// CRYTPO KEY STORE
func (i *Interface) NewUnecryptedKeyStore() crypto.KeyStore {
	return crypto.NewKeyStorePlain(common.DefaultDataDir())
}
func (i *Interface) NewEncryptedKeyStore() crypto.KeyStore {
	return crypto.NewKeyStorePassphrase(common.DefaultDataDir(), crypto.KDFStandard)
}

//func (i *Interface) Peers() []*p2p.Peer      { return i.Peers() }
func (i *Interface) IsListening() bool { return true }

//func (i *Interface) PeerDb() db.Database            { return i.peersDb }
//func (i *Interface) LocalDb() db.Database           { return i.localDb }

// START/QUIT
func (i *Interface) Start() {} // Currently everything starts at initialization
func (i *Interface) Stop() {
	// Stop everything
	i.webUIServer.Stop()
	// Stop p2p networking
	// Stop Tor
	os.Exit(0)
}

// CONFIG
func (i *Interface) GetConfig() *Config {
	return i.config
}
func (i *Interface) SetConfigOption(key string, value string) bool { return false }
func (i *Interface) UnsetConfigOption(key string) bool             { return false }

// WEB UI
func (i *Interface) WebUIStart() bool { return i.webUIServer.Start() }
func (i *Interface) WebUIStop() bool  { return i.webUIServer.Stop() } // Gin is not stoppable

// NETWORK
func (i *Interface) ListPeers() (peers []string)    { return }
func (i *Interface) PeerSuccessor() (peer string)   { return }
func (i *Interface) PeerPredecessor() (peer string) { return }
func (i *Interface) PeerFTable() (peers []string)   { return }
func (i *Interface) NewRing() (ringData string)     { return }
func (i *Interface) ConnectToPeer(peerAddress string) (successful bool) {
	// Eventually get this to return true/false if successful and use this for the return
	go i.p2p.ConnectToPeer(peerAddress, i.tor.ListenPort)
	return true
}
func (in *Interface) RingLookupPeerById(peerId string) (peer string) { return }
func (in *Interface) RingPing(onionAddress string) (pong string)     { return }
func (i *Interface) RingCast(username, body string) (successful bool) {
	message := types.NewMessage(username, body)
	log.Println("[", message.Timestamp, "] ", message.Username, " : ", message.Body)
	i.p2p.Broadcast <- message
	// Eventually get this to return true/false if successful and use this for the return
	return true
}

// DHT
func (ohtInterface *Interface) Put(key string, value string) (successful bool) {
	return
}

func (ohtInterface *Interface) Get(key string) (value string) {
	return
}

func (ohtInterface *Interface) Delete(key string) (value string) {
	return
}
