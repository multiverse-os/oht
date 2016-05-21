package network

import (
	"net/url"
	"time"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Message struct {
	SubProtocol string
	Id          string
	Type        string
	OriginURL   *url.URL
	Timestamp   int64  `json:",omitempty"`
	Username    string `json:",omitempty"`
	Body        string `json:",omitempty"`
}
