package contacts

import (
	"io/ioutil"
	"log"

	"lib/oht/common"
)

type Contacts struct {
	Interface *Interface
}

type Contact struct {
	Id             string `json:"id"`
	LastConnection int64
	OnionHost      string
	Alias          string
	AddRequest     *Request
}

type Request struct {
	Alias   string
	Message string
	Status  int
}

func InitializeContacts() *Contacts {
	if _, err := ioutil.ReadFile(common.AbsolutePath(common.DefaultDataDir(), "contacts.json")); err != nil {
		str := "{}"
		if err = ioutil.WriteFile(common.AbsolutePath(common.DefaultDataDir(), "contacts.json"), []byte(str), 0644); err != nil {
			log.Fatal(err)
		}
	}

	return &Contacts{
		Interface: &Interface{},
	}
}
