package oht

import (
	"./config"
)

type Interface struct {
}

type OHT struct {
	// Channel for shutting down the oht
	shutdownChan chan bool

	// DB interfaces
	// Peer Bolt DB
	// General DHT Bolt DB

	// Handlers
	accountManager *accounts.Manager
	//protocolManager *ProtocolManager -- will this be useful?
	//eventMux *event.TypeMux

	config *Config
}

// msg validator
type Validator interface {
}

type Backend interface {
	AccountManager() *accounts.Manager
	BlockChain() *BlockChain
	TxPool() *TxPool
	ChainDb() ethdb.Database
	DappDb() ethdb.Database
	EventMux() *event.TypeMux
}

// GENERAL INFORMATION
func (interf *Interface) Name() string                      { return s.net.Name }
func (interf *Interface) Version() (version string)         {}
func (interf *Interface) ClientVersion() string             { return s.clientVersion }
func (interf *Interface) Locale() (locale string)           {}
func (interf *Interface) PeerCount() int                    { return s.net.PeerCount() }
func (interf *Interface) MaxPeers() int                     { return s.net.MaxPeers }
func (interf *Interface) Peers() []*p2p.Peer                { return s.net.Peers() }
func (interf *Interface) AccountManager() *accounts.Manager { return s.accountManager }
func (interf *Interface) IsListening() bool                 { return true } // Always listening
func (interf *Interface) PeerDb() ethdb.Database            { return s.dappDb }
func (interf *Interface) LocalDb() ethdb.Database           { return s.dappDb }

// START/QUIT
func (interf *Interface) Start() {}

func (interf *Interface) Quit() {}

// CONFIG
func (interf *Interface) DisplayConfig() (config string, err error) {
	return config.DisplayConfig()
}

func (interf *Interface) SetConfigOption(key string, value string) (successful bool) {

}

func (interf *Interface) UnsetConfigOption(key string) (successful bool) {

}

func (interf *Interface) WebUI(status bool) (successful bool) {

}

// NETWORK
func (interf *Interface) ListPeers() (peers []string) {

}

func (interf *Interface) PeerSuccessor() (peer string) {

}

func (interf *Interface) PeerPredecessor() (peer string) {

}

func (interf *Interface) PeerFTable() (peers []string) {

}

func (interf *Interface) CreateRing() (ringData string) {

}

func (interf *Interface) ConnectToRing(onionAddress string) (successful bool) {

}

func (interf *Interface) RingLookupPeerById(peerId string) (peer string) {

}

func (interf *Interface) RingPing(onionAddress string) (pong string) {

}

func (interf *Interface) RingCast(message string) (successful bool) {

}

// DHT
func (interf *Interface) Put(key string, value string) (successful bool) {

}

func (interf *Interface) Get(key string) (value string) {

}

func (interf *Interface) Delete(key string) (value string) {

}

// ACCOUNT
// This should go into its own folder so it actually is modular
func (interf *Interface) ListAccounts() (accounts []string) {

}

func (interf *Interface) GenerateAccount() (account string) {

}

func (interf *Interface) DeleteAccount(accountId string) (successful bool) {

}

func (interf *Interface) Sign(accountId string) (signature string) {

}

func (interf *Interface) Verify(accountId string, signature string) (successful bool) {

}

func (interf *Interface) Encrypt(accountId string, data string) (encryptedData string) {

}

func (interf *Interface) Decrypt(accountId string, encryptedData string) (data string) {

}

// CONTACTS
func (interf *Interface) ListContacts() (contacts []string) {

}

func (interf *Interface) RequestContact(contactId string, message string) (successful bool) {

}

func (interf *Interface) AddContact(contactId string) (successful bool) {

}

func (interf *Interface) RemoveContact(contactId string) (successful bool) {

}

func (interf *Interface) WhisperToContact(contactId string, message string) (successful bool) {

}

func (interf *Interface) ContactCast(message string) (successful bool) {

}

// CHANNELS
func (interf *Interface) ListChannels() (channels []string) {

}

func (interf *Interface) Channel() (successful bool) {

}

func (interf *Interface) JoinChannel(channelId string) (successful bool) {

}

func (interf *Interface) LeaveChannel(channelId string) (successful bool) {

}

func (interf *Interface) ChannelCast(channelId string, message string) (successful bool) {

}
