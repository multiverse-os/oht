package accounts

import (
	"./../oht/crypto"
)

type Accounts struct {
	manager   *Manager
	Interface *Interface
}

func InitializeAccounts(ks *crypto.KeyStore) *Accounts {
	am := accounts.NewManager(encryptedKeyStore)
	return &Accounts{
		manager:   am,
		Interface: &Interface{Manager: am},
	}
}
