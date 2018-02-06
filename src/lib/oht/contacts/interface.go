package contacts

import ()

type Interface struct {
}

func NewInterface() (i *Interface) {
	return &Interface{}
}

// CONTACTS
func (ohtInterface *Interface) ListContacts() (contacts []string) {
	return
}

func (ohtInterface *Interface) RequestContact(contactId string, message string) (successful bool) {
	return
}

func (ohtInterface *Interface) AddContact(contactId string) (successful bool) {
	return
}

func (ohtInterface *Interface) RemoveContact(contactId string) (successful bool) {
	return
}

func (ohtInterface *Interface) WhisperToContact(contactId string, message string) (successful bool) {
	return
}

func (ohtInterface *Interface) ContactCast(message string) (successful bool) {
	return
}
