package accounts

import (
	"../oht"
)

type Accounts struct {
	oht       *OHT
	manager   *Manager
	Interface *Interface
}

func (oht *OHT) InitializeAccounts(ks *crypto.KeyStore) *Accounts {
	am := Manager(encryptedKeyStore)
	return &Accounts{
		oht:       oht,
		manager:   am,
		Interface: &Interface{Manager: am},
	}
}
