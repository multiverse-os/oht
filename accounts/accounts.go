package accounts

import ()

type Accounts struct {
	Manager   *Manager
	Interface *Interface
}

func InitializeAccounts() *Accounts {
	return &Accounts{
		Manager:   &Manager{},
		Interface: &Interface{},
	}
}
