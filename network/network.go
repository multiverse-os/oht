package network

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"net/http"
	"time"
)

type EventFunc func(Manager *NetworkManager, Peer *Peer)

type NetworkManager struct {
	PrivateKey *ecdsa.PrivateKey
	Server     *Server

	MaxPeers        int
	MaxPendingPeers int

	Peers      map[*Peer]bool
	Broadcast  chan Message
	Receive    chan Message
	Register   chan *Peer
	Unregister chan *Peer
	OnConnect  EventFunc
	OnClose    EventFunc

	lastLookup time.Time
}

var Manager = NetworkManager{
	// Need to add a new message with custom type to send out
	// Info about current node, like onion address
	Server:          &Server{},
	MaxPeers:        8,
	MaxPendingPeers: 8,
	Broadcast:       make(chan Message, maxMessageSize),
	Receive:         make(chan Message, maxMessageSize),
	Register:        make(chan *Peer, maxMessageSize),
	Unregister:      make(chan *Peer, maxMessageSize),
	Peers:           make(map[*Peer]bool, maxMessageSize),
	OnConnect:       nil,
	OnClose:         nil,
	lastLookup:      time.Now(),
}

func (Manager *NetworkManager) Start(port int) {
	Manager.Server.Initialize(port)
	for {
		select {
		case p := <-Manager.Register:
			Manager.Peers[p] = true
			log.Println("Peer connection established: ", p.OnionHost)
			fmt.Printf("whisper> ")
			if Manager.OnConnect != nil {
				go Manager.OnConnect(Manager, p)
			}
		case p := <-Manager.Unregister:
			if _, ok := Manager.Peers[p]; ok {
				delete(Manager.Peers, p)
				close(p.Send)
				if Manager.OnClose != nil {
					go Manager.OnClose(Manager, p)
				}
				p.Websocket.Close()
			}
		case m := <-Manager.Broadcast:
			for p := range Manager.Peers {
				select {
				case p.Send <- m:
				default:
					close(p.Send)
					delete(Manager.Peers, p)
				}
			}
		case m := <-Manager.Receive:
			log.Println("")
			log.Println("[", m.Timestamp, "] ", m.Username, " : ", m.Body)
			fmt.Printf("whisper> ")
		}
	}

}

func (Manager *NetworkManager) DumpPeers() {
	for p := range Manager.Peers {
		log.Println("Connection")
		log.Println("Connection: ", p.OnionHost)
	}
}

// Serve handles websocket requests from the peer
func (Manager *NetworkManager) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	p := &Peer{Send: make(chan Message, 256), Websocket: ws}
	Manager.Register <- p
	go p.writeMessages()
	p.readMessages()
}
