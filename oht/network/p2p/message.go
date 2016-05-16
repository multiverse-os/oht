package network

import (
	"time"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Message struct {
	Type       string
	Id         string
	Timestamp  int64  `json:",omitempty"`
	OriginHost string `json:",omitempty"`
	Username   string `json:",omitempty"`
	Body       string `json:",omitempty"`
}
