package oht

import (
	"log"

	"./common"
	"./crypto"
	"./network"
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

// STATE INFORMATION
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

// KEY STORE
func (i *Interface) NewUnecryptedKeyStore() crypto.KeyStore {
	return crypto.NewKeyStorePlain(common.DefaultDataDir() + "/keys")
}
func (i *Interface) NewEncryptedKeyStore() crypto.KeyStore {
	return crypto.NewKeyStorePassphrase(common.DefaultDataDir()+"/keys", crypto.KDFStandard)
}

// CONFIG INTERFACE
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

// TOR INTERFACE
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
	i.tor.Stop()
	i.tor.DeleteOnionFiles()
	return i.tor.Start()
}

// NETWORK INTERFACE
func (i *Interface) ListPeers() (peers []string)    { return }
func (i *Interface) PeerSuccessor() (peer string)   { return }
func (i *Interface) PeerPredecessor() (peer string) { return }
func (i *Interface) PeerFTable() (peers []string)   { return }
func (i *Interface) NewRing() (ringData string)     { return }
func (i *Interface) ConnectToPeer(peerAddress string) bool {
	if i.tor.Online {
		go i.p2p.ConnectToPeer(peerAddress, i.tor.ListenPort)
		return true
	} else {
		return false
	}
}
func (in *Interface) RingLookupPeerById(peerId string) string { return "" }
func (in *Interface) RingPing(onionAddress string) bool       { return false }
func (i *Interface) RingCast(username, body string) bool {
	message := types.NewMessage(username, body)
	if body != "" {
		log.Println("[", message.Timestamp, "] ", message.Username, " : ", message.Body)
		i.p2p.Broadcast <- message
	}
	return true
}

// DHT INTERFACE
func (ohtInterface *Interface) Put(key string, value string) (successful bool) {
	return
}
func (ohtInterface *Interface) Get(key string) (value string) {
	return
}
func (ohtInterface *Interface) Delete(key string) (value string) {
	return
}

// WEB UI INTERFACE
func (i *Interface) WebUIOnline() bool {
	return i.webUI.Server.Online
}
func (i *Interface) WebUIStart() bool {
	if i.tor.Online {
		err := i.webUI.Server.Start()
		return (err == nil)
	} else {
		return false
	}
}
func (i *Interface) WebUIStop() bool {
	return i.webUI.Server.Stop()
}
