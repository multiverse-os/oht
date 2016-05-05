package oht

import (
	"os"

	"./../accounts"
	"./network"
)

type Interface struct {
	config *Config
	tor    *network.TorProcess
}

func NewInterface(c *Config, t *network.TorProcess) (i *Interface) {
	return &Interface{
		config: c,
		tor:    t,
	}
}

// msg validator
type Validator interface {
}

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

//func (in *Interface) Peers() []*p2p.Peer      { return in.Peers() }
func (othInterface *Interface) IsListening() bool { return true }

//func (in *Interface) PeerDb() db.Database            { return in.peersDb }
//func (in *Interface) LocalDb() db.Database           { return in.localDb }

// START/QUIT
func (othInterface *Interface) Start() {}
func (othInterface *Interface) Quit()  { os.Exit(0) }

// CONFIG
func (in *Interface) GetConfig() (config []byte, err error) { return in.config.getConfig() }

func (in *Interface) SetConfigOption(key string, value string) (successful bool) {
	return
}

func (in *Interface) UnsetConfigOption(key string) (successful bool) {
	return
}

func (in *Interface) WebUI(status bool) (successful bool) {
	if status == 1 {
		webui.InitializeServer(oht.Interface.TorWebUIOnionHost(), oht.Interface.TorWebUIPort())
		log.Printf("\nWeb UI :  " + oht.Interface.TorWebUIOnionHost() + ":" + strconv.Itoa(oht.Interface.TorWebUIPort()))
	} else if status == 2 {

	}
	return
}

// NETWORK
func (in *Interface) ListPeers() (peers []string) {
	return
}

func (in *Interface) PeerSuccessor() (peer string) {
	return
}

func (in *Interface) PeerPredecessor() (peer string) {
	return
}

func (in *Interface) PeerFTable() (peers []string) {
	return
}

func (in *Interface) CreateRing() (ringData string) {
	return
}

func (in *Interface) ConnectToRing(onionAddress string) (successful bool) {
	return
}

func (in *Interface) RingLookupPeerById(peerId string) (peer string) {
	return
}

func (in *Interface) RingPing(onionAddress string) (pong string) {
	return
}

func (in *Interface) RingCast(message string) (successful bool) {
	return
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
