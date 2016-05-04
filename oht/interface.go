package oht

import ()

type Interface struct {
}

// GENERAL INFORMATION
func (interf *Interface) Version() (version string) {
}

// CONFIG
func (interf *Interface) DisplayConfig() (config string) {

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

// QUIT
func (interf *Interface) Quit() {

}
