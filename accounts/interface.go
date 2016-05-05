package accounts

import ()

type Interface struct {
	Manager *Manager
}

func NewInterface(am *Manager) (i *Interface) {
	return &Interface{
		Manager: am,
	}
}

func (i *Interface) AccountManager() *Manager { return i.Manager }

// ACCOUNT
func (i *Interface) ListAccounts() (accounts []string) {
	return
}

func (i *Interface) GenerateAccount() (account string) {
	return
}

func (i *Interface) DeleteAccount(accountId string) (successful bool) {
	return
}

func (i *Interface) Sign(accountId string) (signature string) {
	return
}

func (i *Interface) Verify(accountId string, signature string) (successful bool) {
	return
}

func (i *Interface) Encrypt(accountId string, data string) (encryptedData string) {
	return
}

func (i *Interface) Decrypt(accountId string, encryptedData string) (data string) {
	return
}
