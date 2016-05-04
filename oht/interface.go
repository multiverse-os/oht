package oht

import (
	"./config"
)

type OHTInterface struct {
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

// GENERAL INFORMATION
func (ohtInterface *OHTInterface) Name() string                      { return s.net.Name }
func (ohtInterface *OHTInterface) Version() (version string)         {}
func (ohtInterface *OHTInterface) ClientVersion() string             { return s.clientVersion }
func (ohtInterface *OHTInterface) Locale() (locale string)           {}
func (ohtInterface *OHTInterface) PeerCount() int                    { return s.net.PeerCount() }
func (ohtInterface *OHTInterface) MaxPeers() int                     { return s.net.MaxPeers }
func (ohtInterface *OHTInterface) Peers() []*p2p.Peer                { return s.net.Peers() }
func (ohtInterface *OHTInterface) AccountManager() *accounts.Manager { return s.accountManager }
func (ohtInterface *OHTInterface) IsListening() bool                 { return true } // Always listening
func (ohtInterface *OHTInterface) PeerDb() ethdb.Database            { return s.dappDb }
func (ohtInterface *OHTInterface) LocalDb() ethdb.Database           { return s.dappDb }

// START/QUIT
func (ohtInterface *OHTInterface) Start() {}
func (ohtInterface *OHTInterface) Quit()  {}

// CONFIG
func (ohtInterface *OHTInterface) DisplayConfig() (config string, err error) {
	return config.DisplayConfig()
}

func (ohtInterface *OHTInterface) SetConfigOption(key string, value string) (successful bool) {

}

func (ohtInterface *OHTInterface) UnsetConfigOption(key string) (successful bool) {

}

func (OhtInterface *OHTInterface) WebUI(status bool) (successful bool) {

}

// NETWORK
func (ohtInterface *OHTInterface) ListPeers() (peers []string) {

}

func (ohtInterface *OHTInterface) PeerSuccessor() (peer string) {

}

func (ohtInterface *OHTInterface) PeerPredecessor() (peer string) {

}

func (ohtInterface *OHTInterface) PeerFTable() (peers []string) {

}

func (ohtInterface *OHTInterface) CreateRing() (ringData string) {

}

func (ohtInterface *OHTInterface) ConnectToRing(onionAddress string) (successful bool) {

}

func (ohtInterface *OHTInterface) RingLookupPeerById(peerId string) (peer string) {

}

func (ohtInterface *OHTInterface) RingPing(onionAddress string) (pong string) {

}

func (ohtInterface *OHTInterface) RingCast(message string) (successful bool) {

}

// DHT
func (ohtInterface *OHTInterface) Put(key string, value string) (successful bool) {

}

func (ohtInterface *OHTInterface) Get(key string) (value string) {

}

func (ohtInterface *OHTInterface) Delete(key string) (value string) {

}

// ACCOUNT
// This should go into its own folder so it actually is modular
func (ohtInterface *OHTInterface) ListAccounts() (accounts []string) {

}

func (ohtInterface *OHTInterface) GenerateAccount() (account string) {

}

func (ohtInterface *OHTInterface) DeleteAccount(accountId string) (successful bool) {

}

func (ohtInterface *OHTInterface) Sign(accountId string) (signature string) {

}

func (ohtInterface *OHTInterface) Verify(accountId string, signature string) (successful bool) {

}

func (ohtInterface *OHTInterface) Encrypt(accountId string, data string) (encryptedData string) {

}

func (ohtInterface *OHTInterface) Decrypt(accountId string, encryptedData string) (data string) {

}

// CONTACTS
func (ohtInterface *OHTInterface) ListContacts() (contacts []string) {

}

func (ohtInterface *OHTInterface) RequestContact(contactId string, message string) (successful bool) {

}

func (ohtInterface *OHTInterface) AddContact(contactId string) (successful bool) {

}

func (ohtInterface *OHTInterface) RemoveContact(contactId string) (successful bool) {

}

func (ohtInterface *OHTInterface) WhisperToContact(contactId string, message string) (successful bool) {

}

func (ohtInterface *OHTInterface) ContactCast(message string) (successful bool) {

}

// CHANNELS
func (ohtInterface *OHTInterface) ListChannels() (channels []string) {

}

func (ohtInterface *OHTInterface) Channel() (successful bool) {

}

func (ohtInterface *OHTInterface) JoinChannel(channelId string) (successful bool) {

}

func (ohtInterface *OHTInterface) LeaveChannel(channelId string) (successful bool) {

}

func (ohtInterface *OHTInterface) ChannelCast(channelId string, message string) (successful bool) {

}
