package types

import (
	"time"

	"github.com/pborman/uuid"
)

type Message struct {
	Type      string
	Id        string
	Timestamp int64  `json:",omitempty"`
	Username  string `json:",omitempty"`
	Body      string `json:",omitempty"`
}

func NewMessage(username, body string) Message {
	return Message{
		Id:        uuid.New(),
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Username:  username,
		Body:      body,
	}
}
