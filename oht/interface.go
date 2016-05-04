package oht

import ()

type Interface struct {
}

// Application Information
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
func (interf *Interface) ListPeers(peers []string) {

}

func (interf *Interface) PeerSuccessor(peer string) {

}

func (interf *Interface) PeerPredecessor(peer string) {

}

func (interf *Interface) PeerFTable(peers []string) {

}

func (interf *Interface) CreateRing(successful bool) {

}

func (interf *Interface) ConnectToRing(onionAddress string) (successful bool) {

}

func (interf *Interface) RingLookupPeerById(id string) (peer string) {

}
