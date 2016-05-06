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
	config *Config
	tor    *network.TorProcess
	webUI  *webui.WebUI
	p2p    *p2p.Manager
}

func NewInterface(c *Config, t *network.TorProcess, w *webui.WebUI, p *p2p.Manager) (i *Interface) {
	return &Interface{
		config: c,
		tor:    t,
		webUI:  w,
		p2p:    p,
	}
}

// GENERAL INFORMATION
func (i *Interface) ClientName() string    { return i.config.ClientName }
func (i *Interface) ClientVersion() string { return i.config.clientVersion() }
func (i *Interface) ClientInfo() string    { return i.config.clientInfo() }
func (i *Interface) Locale() string        { return i.config.Locale }

func (i *Interface) TorListenPort() string     { return i.tor.ListenPort }
func (i *Interface) TorSocksPort() string      { return i.tor.SocksPort }
func (i *Interface) TorControlPort() string    { return i.tor.ControlPort }
func (i *Interface) TorWebUIPort() string      { return i.tor.WebUIPort }
func (i *Interface) TorOnionHost() string      { return i.tor.OnionHost }
func (i *Interface) TorWebUIOnionHost() string { return i.tor.WebUIOnionHost }

func (i *Interface) MaxPendingPeers() int { return i.config.MaxPendingPeers }
func (i *Interface) MaxPeers() int        { return i.config.MaxPeers }

// QUIT
func (i *Interface) Stop() {
	if i.webUI.Server.Online {
		i.webUI.Server.Stop()
	}
	// Stop everything
	// Stop p2p networking
	// Stop Tor
	os.Exit(0)
}

// KEY STORE
func (i *Interface) NewUnecryptedKeyStore() crypto.KeyStore {
	return crypto.NewKeyStorePlain(common.DefaultDataDir() + "/keys")
}
func (i *Interface) NewEncryptedKeyStore() crypto.KeyStore {
	return crypto.NewKeyStorePassphrase(common.DefaultDataDir()+"/keys", crypto.KDFStandard)
}

// CONFIG
func (i *Interface) Config() *Config {
	return i.config
}
func (i *Interface) ConfigSetOption(key, value string) bool {
	return i.config.setConfigOption(key, value)
}
func (i *Interface) ConfigUnsetOption(key string) bool {
	return i.config.unsetConfigOption(key)
}
func (i *Interface) ConfigSave() bool {
	return i.config.saveConfiguration()
}

// TOR
func (i *Interface) TorOnline() bool {
	return i.tor.Online
}
func (i *Interface) TorStart() bool {
	return i.tor.Start()
}
func (i *Interface) TorStop() bool {
	return i.tor.Stop()
}
func (i *Interface) TorCycleIdentity() bool {
	i.tor.Cycle()
	return true
}
func (i *Interface) TorCycleOnionAddresses() bool {
	return true
}

// NETWORK
func (i *Interface) ListPeers() (peers []string)    { return }
func (i *Interface) PeerSuccessor() (peer string)   { return }
func (i *Interface) PeerPredecessor() (peer string) { return }
func (i *Interface) PeerFTable() (peers []string)   { return }
func (i *Interface) NewRing() (ringData string)     { return }
func (i *Interface) ConnectToPeer(peerAddress string) (successful bool) {
	go i.p2p.ConnectToPeer(peerAddress, i.tor.ListenPort)
	return true
}
func (in *Interface) RingLookupPeerById(peerId string) (peer string) { return }
func (in *Interface) RingPing(onionAddress string) (pong string)     { return }
func (i *Interface) RingCast(username, body string) (successful bool) {
	message := types.NewMessage(username, body)
	log.Println("[", message.Timestamp, "] ", message.Username, " : ", message.Body)
	i.p2p.Broadcast <- message
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

// WEB UI
func (i *Interface) WebUIOnline() bool {
	return i.webUI.Server.Online
}
func (i *Interface) WebUIStart() bool {
	err := i.webUI.Server.Start()
	return (err == nil)
}
func (i *Interface) WebUIStop() bool {
	return i.webUI.Server.Stop()
}
