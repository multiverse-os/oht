package accounts

import ()

type Interface struct {
	accountManager *accounts.Manager
}

func NewInterface(am *accounts.Manager) (i *Interface) {
	return &Interface{
		accountManager: am,
	}
}

func (i *Interface) AccountManager() *accounts.Manager { return othInterface.accountManager }

// DEV

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
