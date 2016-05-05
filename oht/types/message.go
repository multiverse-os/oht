package types

import (
	"time"

	"../network"

	"github.com/pborman/uuid"
)

type Message struct {
	Type      string
	Id        string
	Timestamp *time.Time `json:",omitempty"`
	Username  string     `json:",omitempty"`
	Body      string     `json:",omitempty"`
}

func NewMessage(username, body string) (message *Message) {
	return &Message{
		Id:        uuid.New(),
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Username:  username,
		Body:      body,
	}
}

// Will need to add Ringcast and ways to narrow the broadcast range
func (message *Message) Broadcast() {
	network.Manager.Broadcast <- message
}
