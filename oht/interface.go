package oht

import (
	"os"

	"./../accounts"
	"./network"
	"./types"
)

type Interface struct {
	config      *Config
	tor         *network.TorProcess
	webUIServer *webui.Server
}

func NewInterface(c *Config, t *network.TorProcess, w *webui.Server) (i *Interface) {
	return &Interface{
		config:      c,
		tor:         t,
		webUIServer: w,
	}
}

// msg validator
type Validator interface{}

// GENERAL INFORMATION
func (othInterface *Interface) ClientName() string    { return othInterface.config.clientName }
func (othInterface *Interface) ClientVersion() string { return othInterface.config.clientVersion() }
func (othInterface *Interface) ClientInfo() string    { return othInterface.config.clientInfo() }

// TOR INFORMATION
func (othInterface *Interface) TorListenPort() int        { return othInterface.tor.ListenPort }
func (othInterface *Interface) TorSocksPort() int         { return othInterface.tor.SocksPort }
func (othInterface *Interface) TorControlPort() int       { return othInterface.tor.ControlPort }
func (othInterface *Interface) TorWebUIPort() int         { return othInterface.tor.WebUIPort }
func (othInterface *Interface) TorOnionHost() string      { return othInterface.tor.OnionHost }
func (othInterface *Interface) TorWebUIOnionHost() string { return othInterface.tor.WebUIOnionHost }

func (othInterface *Interface) ProtocolVersion() {}
func (othInterface *Interface) Locale() string   { return "en" }
func (othInterface *Interface) PeerCount() int   { return othInterface.PeerCount() }
func (othInterface *Interface) MaxPeers() int    { return othInterface.MaxPeers() }

// CRYTPO KEY STORE
func (i *Interface) GenerateUnecryptedKeystore() *crypto.KeyStore {
	return crypto.NewKeyStorePlain(common.DefaultDataDir())
}

func (i *Interface) GenerateEncryptedKeystore(password string) *crypto.KeyStore {
	return crypto.NewKeyStorePassphrase(common.DefaultDataDir(), crypto.KDFStandard)
}

//func (i *Interface) Peers() []*p2p.Peer      { return i.Peers() }
func (i *Interface) IsListening() bool { return true }

//func (i *Interface) PeerDb() db.Database            { return i.peersDb }
//func (i *Interface) LocalDb() db.Database           { return i.localDb }

// START/QUIT
func (i *Interface) Start() {}
func (i *Interface) Quit() {
	// Stop everything
	i.webUIServer.Stop()
	// Stop p2p networking
	// Stop Tor
	os.Exit(0)
}

// CONFIG
func (i *Interface) GetConfig() (config []byte, err error) { return i.config.getConfig() }

func (i *Interface) SetConfigOption(key string, value string) bool {
	return
}

func (i *Interface) UnsetConfigOption(key string) bool {
	return
}

// WEB UI
func (i *Interface) WebUIStart() bool {
	return i.webUIServer.Start()
}

func (i *Interface) WebUIStop() bool {
	return i.webUIServer.Stop()
}

// NETWORK
func (i *Interface) ListPeers() (peers []string) {
	return
}

func (i *Interface) PeerSuccessor() (peer string) {
	return
}

func (i *Interface) PeerPredecessor() (peer string) {
	return
}

func (i *Interface) PeerFTable() (peers []string) {
	return
}

func (i *Interface) CreateRing() (ringData string) {
	return
}

func (i *Interface) ConnectToPeer(peerAddress string) (successful bool) {
	// Eventually get this to return true/false if successful and use this for the return
	go network.ConnectToPeer(peerAddress)
	return true
}

func (in *Interface) RingLookupPeerById(peerId string) (peer string) {
	return
}

func (in *Interface) RingPing(onionAddress string) (pong string) {
	return
}

func (i *Interface) RingCast(username, message string) (successful bool) {
	message := types.Message.NewMessage(username, message)
	// Eventually get this to return true/false if successful and use this for the return
	message.Broadcast()
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
