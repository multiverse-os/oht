package oht

import (
	"os"

	"./../accounts"
	"./../contacts"
	"./common"
)

type OHTInterface struct {
	clientName   string
	versionMajor int
	versionMinor int
	versionPatch int

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

func NewOHTInterface() (ohtInterface *OHTInterface) {
	InitializeConfig()
	contacts.InitializeContacts()
	// Initialize Data Directory
	if !common.FileExist(common.DefaultDataDir()) {
		os.MkdirAll(common.DefaultDataDir(), os.ModePerm)
	}
	return &OHTInterface{
		clientName:   "oth",
		versionMajor: 0,
		versionMinor: 1,
		versionPatch: 0,
		config:       &Config{},
	}
}

// msg validator
type Validator interface {
}

// GENERAL INFORMATION
func (ohtInterface *OHTInterface) ClientName() string    { return ohtInterface.config.ClientName }
func (ohtInterface *OHTInterface) ClientVersion() string { return ohtInterface.config.ClientVersion() }
func (ohtInterface *OHTInterface) ProtocolVersion()      {}
func (ohtInterface *OHTInterface) Locale() string        { return "en" }
func (ohtInterface *OHTInterface) PeerCount() int        { return ohtInterface.PeerCount() }
func (ohtInterface *OHTInterface) MaxPeers() int         { return ohtInterface.MaxPeers() }

//func (ohtInterface *OHTInterface) Peers() []*p2p.Peer      { return ohtInterface.Peers() }
func (ohtInterface *OHTInterface) AccountManager() *accounts.Manager {
	return ohtInterface.accountManager
}
func (ohtInterface *OHTInterface) IsListening() bool { return true } // Always listening
//func (ohtInterface *OHTInterface) PeerDb() ethdb.Database            { return ohtInterface.peersDb }
//func (ohtInterface *OHTInterface) LocalDb() ethdb.Database           { return ohtInterface.localDb }

// START/QUIT
func (ohtInterface *OHTInterface) Start() {}
func (ohtInterface *OHTInterface) Quit()  { os.Exit(0) }

// CONFIG
func (ohtInterface *OHTInterface) GetConfig() (config []byte, err error) {
	return ohtInterface.config.GetConfig()
}

func (ohtInterface *OHTInterface) SetConfigOption(key string, value string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) UnsetConfigOption(key string) (successful bool) {
	return
}

func (OhtInterface *OHTInterface) WebUI(status bool) (successful bool) {
	return
}

// NETWORK
func (ohtInterface *OHTInterface) ListPeers() (peers []string) {
	return
}

func (ohtInterface *OHTInterface) PeerSuccessor() (peer string) {
	return
}

func (ohtInterface *OHTInterface) PeerPredecessor() (peer string) {
	return
}

func (ohtInterface *OHTInterface) PeerFTable() (peers []string) {
	return
}

func (ohtInterface *OHTInterface) CreateRing() (ringData string) {
	return
}

func (ohtInterface *OHTInterface) ConnectToRing(onionAddress string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) RingLookupPeerById(peerId string) (peer string) {
	return
}

func (ohtInterface *OHTInterface) RingPing(onionAddress string) (pong string) {
	return
}

func (ohtInterface *OHTInterface) RingCast(message string) (successful bool) {
	return
}

// DHT
func (ohtInterface *OHTInterface) Put(key string, value string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) Get(key string) (value string) {
	return
}

func (ohtInterface *OHTInterface) Delete(key string) (value string) {
	return
}

// ACCOUNT
// This should go into its own folder so it actually is modular
func (ohtInterface *OHTInterface) ListAccounts() (accounts []string) {
	return
}

func (ohtInterface *OHTInterface) GenerateAccount() (account string) {
	return
}

func (ohtInterface *OHTInterface) DeleteAccount(accountId string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) Sign(accountId string) (signature string) {
	return
}

func (ohtInterface *OHTInterface) Verify(accountId string, signature string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) Encrypt(accountId string, data string) (encryptedData string) {
	return
}

func (ohtInterface *OHTInterface) Decrypt(accountId string, encryptedData string) (data string) {
	return
}

// CONTACTS
func (ohtInterface *OHTInterface) ListContacts() (contacts []string) {
	return
}

func (ohtInterface *OHTInterface) RequestContact(contactId string, message string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) AddContact(contactId string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) RemoveContact(contactId string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) WhisperToContact(contactId string, message string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) ContactCast(message string) (successful bool) {
	return
}

// CHANNELS
func (ohtInterface *OHTInterface) ListChannels() (channels []string) {
	return
}

func (ohtInterface *OHTInterface) Channel() (successful bool) {
	return
}

func (ohtInterface *OHTInterface) JoinChannel(channelId string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) LeaveChannel(channelId string) (successful bool) {
	return
}

func (ohtInterface *OHTInterface) ChannelCast(channelId string, message string) (successful bool) {
	return
}
